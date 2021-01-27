/*
Zahhak2, a Golang console game.
Copyright (C) 2021 Aryo Pehlewan aryopehlewan@hotmail.com
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
	"time"

	"github.com/gorilla/websocket"

	z "../common"
)

type Connection struct {
	address string
	hub     *BroadcastHub
	ws      *websocket.Conn
	send    chan []byte
}

func (c *Connection) readPump() {
	defer c.recover()

	defer func() {
		log.Println("Connection: " + c.address + " readPump(): ws.Close()")
		c.hub.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.ws.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				z.LogError(errors.New("Connection: " + c.address + " " + " websocket.IsUnexpectedCloseError"))
			} else {
				z.LogError(errors.New("Connection: " + c.address + " " + err.Error()))
			}

			break
		}

		c.hub.incoming <- message
	}
}

func (c *Connection) writePump() {
	defer c.recover()

	ticker := time.NewTicker(pingPeriod)

	defer func() {
		log.Println("Connection: " + c.address + " writePump(): ws.Close()")
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				z.LogError(errors.New("Connection: " + c.address + " " + " websocket.CloseMessage"))
				c.write(websocket.CloseMessage, []byte{})

				return
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				z.LogError(errors.New("Connection: " + c.address + " " + err.Error()))

				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				z.LogError(errors.New("Connection: " + c.address + " " + err.Error()))

				return
			}
		}
	}
}

func (c *Connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))

	return c.ws.WriteMessage(mt, payload)
}

func (c *Connection) recover() {
	if r := recover(); r != nil {
		e, ok := r.(error)

		if !ok {
			e = fmt.Errorf("%v", r)
			e = errors.New("Connection: " + c.address + " " + e.Error())
			z.LogPanic(e)
		} else {
			e = errors.New("Connection: " + c.address + " " + e.Error())
			z.LogPanic(e)
		}
	}
}
