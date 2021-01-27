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
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	z "../common"
)

type Server struct {
	address      string
	hub          *BroadcastHub
	eventManager *z.EventManager
}

func NewServer(eventManager *z.EventManager, port string, incoming, outgoing chan []byte) *Server {
	address := "127.0.0.1:" + port

	hub := NewBroadcastHub(incoming, outgoing)

	return &Server{
		address:      address,
		hub:          hub,
		eventManager: eventManager,
	}
}

func (s *Server) Run() {
	s.hub.Run()

	go func() {
		http.HandleFunc("/", s.handler)
		e := http.ListenAndServe(s.address, nil)

		if e != nil {
			panic(e)
		}
	}()
}

func (s *Server) Count() int {
	return s.hub.Count()
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	address := r.RemoteAddr

	defer s.recover(address)

	defer log.Println("Server: " + address + " Closed websocket")

	log.Println("Server: " + address + " New HTTP connection")

	upgrader := websocket.Upgrader{
		ReadBufferSize:  65335,
		WriteBufferSize: 65335,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	log.Println("Server: " + address + " Attempting to upgrade to websocket")

	ws, e := upgrader.Upgrade(w, r, nil)

	if e != nil {
		z.LogError(errors.New("Server: " + address + " " + e.Error()))

		return
	}

	log.Println("Server: " + address + " Calling to NewClient")

	p, er := s.eventManager.Fire("NewClient")

	if er != nil {
		z.LogError(errors.New("Server: " + address + " " + er.Error()))

		return
	}

	log.Println("Server: " + address + " Returning from NewClient")

	v := p[0]
	json := v.String()

	bs := []byte(json)

	log.Println("Server: " + address + " Writing to websocket")

	if err := s.write(ws, websocket.TextMessage, bs); err != nil {
		z.LogError(errors.New("Server: " + address + " " + err.Error()))

		return
	}

	log.Println("Server: " + address + " Creating Connection")

	c := &Connection{address: address, hub: s.hub, send: make(chan []byte, 1024), ws: ws}

	log.Println("Server: " + address + " Registering Connection")

	s.hub.register <- c

	log.Println("Server: " + address + " Starting pumps")

	go c.writePump()
	c.readPump()
}

func (s *Server) write(ws *websocket.Conn, mt int, payload []byte) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))

	return ws.WriteMessage(mt, payload)
}

func (s *Server) recover(address string) {
	if r := recover(); r != nil {
		e, ok := r.(error)

		if !ok {
			e = fmt.Errorf("%v", r)
			e = errors.New("Server: " + address + " " + e.Error())
			z.LogPanic(e)
		} else {
			e = errors.New("Server: " + address + " " + e.Error())
			z.LogPanic(e)
		}
	}
}
