/*
Zahhak2, a Golang multiplayer console game.
Copyright (C) 2016 Aryo Pehlewan aryopehlewan@hotmail.com
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package networking

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	z "../common"
)

type Client struct {
	address    string
	connection *websocket.Conn
	incoming   chan []byte
	outgoing   chan []byte
}

func NewClient(host, port string, incoming, outgoing chan []byte) *Client {
	address := host + ":" + port

	return &Client{
		address:  address,
		incoming: incoming,
		outgoing: outgoing,
	}
}

func (c *Client) Run() {
	defer c.recover()

	var e error

	u := url.URL{Scheme: "ws", Host: c.address, Path: "/"}

	c.connection, _, e = websocket.DefaultDialer.Dial(u.String(), nil)

	if e != nil {
		z.LogError(errors.New("Client: " + e.Error()))
	}

	go c.writePump()
	go c.readPump()
}

func (c *Client) readPump() {
	defer c.recover()

	defer func() {
		log.Printf("Client: readPump(): Close()")

		c.connection.Close()
	}()

	c.connection.SetPongHandler(func(string) error { c.connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				z.LogError(errors.New("Client: readPump(): IsUnexpectedCloseError"))
			} else {
				z.LogError(errors.New("Client: readPump(): " + err.Error()))
			}

			break
		}

		c.incoming <- message
	}
}

func (c *Client) writePump() {
	defer c.recover()

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.Printf("Client: writePump(): Close()")
		ticker.Stop()
		c.connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.outgoing:
			if !ok {
				z.LogError(errors.New("Client: writePump(): CloseMessage"))
				c.write(websocket.CloseMessage, []byte{})

				return
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				z.LogError(errors.New("Client: writePump(): " + err.Error()))

				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				z.LogError(errors.New("Client: writePump(): " + err.Error()))

				return
			}
		}
	}
}

func (c *Client) write(mt int, payload []byte) error {
	c.connection.SetWriteDeadline(time.Now().Add(writeWait))

	return c.connection.WriteMessage(mt, payload)
}

func (c *Client) Count() int {
	if c.connection == nil {
		return 0
	}

	return 1
}

func (c *Client) recover() {
	if r := recover(); r != nil {
		e, ok := r.(error)

		if !ok {
			e = fmt.Errorf("%v", r)
			e = errors.New("Client: " + e.Error())
			z.LogPanic(e)
		} else {
			e = errors.New("Client: " + e.Error())
			z.LogPanic(e)
		}
	}
}
