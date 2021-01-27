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

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	z "./common"

	zn "./networking"
)

type NetworkManager struct {
	gameID  string
	session string
	mode    string

	server     bool
	incoming   chan []byte
	outgoing   chan []byte
	connection z.IConnection

	started  bool
	Messages chan *z.Message
}

func NewNetworkManager(gameID, session string, em *z.EventManager, host, port string, server bool) *NetworkManager {
	in := make(chan []byte, 1024)
	out := make(chan []byte, 1024)
	ms := make(chan *z.Message, 1024)

	var c z.IConnection
	mode := "NetworkManager: "

	if server {
		c = zn.NewServer(em, port, in, out)
		mode = "Server: "
	} else {
		c = zn.NewClient(host, port, in, out)
		mode = "Client: "
	}

	return &NetworkManager{
		gameID:     gameID,
		session:    session,
		mode:       mode,
		server:     server,
		incoming:   in,
		outgoing:   out,
		connection: c,
		Messages:   ms,
	}
}

func (nm *NetworkManager) Run() {
	nm.connection.Run()
	nm.started = true

	go func() {
		defer nm.recover()

		for {
			var m *z.Message
			var e error

			bs, ok := <-nm.incoming

			if !ok {
				e = errors.New(nm.mode + "Incoming: channel error")
				z.LogError(e)
				m = nm.error(e)
			} else {
				m, e = nm.Receive(bs)

				if e != nil {
					e = errors.New(nm.mode + "Incoming: " + e.Error())
					z.LogError(e)
					m = nm.error(e)
				}
			}

			if m != nil {
				nm.Messages <- m
			}
		}
	}()
}

func (nm *NetworkManager) Receive(bs []byte) (*z.Message, error) {
	m := &z.Message{}

	e := json.Unmarshal(bs, m)

	return m, e
}

func (nm *NetworkManager) Send(m *z.Message) error {
	defer nm.recover()

	var e error

	if !nm.started {
		e = errors.New("NetworkManager not started")

		return e
	}

	if nm.connection.Count() == 0 {
		e = errors.New("No connection")

		return e
	}

	bs, e := json.Marshal(m)

	if e != nil {
		return e
	}

	select {
	case nm.outgoing <- bs:

	default:
	}

	return e
}

func (nm *NetworkManager) error(e error) *z.Message {
	m := &z.Message{}
	m.Class = "NetworkManager"
	m.ID = "0"
	m.Action = "Error"
	m.Params["Session"] = nm.session
	m.Params["GameID"] = nm.gameID
	m.Params["Exception"] = e.Error()

	return m
}

func (nm *NetworkManager) recover() {
	if r := recover(); r != nil {
		e, ok := r.(error)

		if !ok {
			e = fmt.Errorf("%v", r)
			e = errors.New(nm.mode + "Recover(): " + e.Error())
			z.LogPanic(e)
		} else {
			e = errors.New(nm.mode + "Recover(): " + e.Error())
			z.LogPanic(e)
		}
	}
}
