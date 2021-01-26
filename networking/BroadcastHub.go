/*
Zahhak2, a Golang multiplayer console game.
Copyright (C) 2016 Aryo Pehlewan feylikurds@gmail.com
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

type BroadcastHub struct {
	connections map[*Connection]bool

	register   chan *Connection
	unregister chan *Connection

	incoming chan []byte
	outgoing chan []byte
}

func NewBroadcastHub(incoming chan []byte, outgoing chan []byte) *BroadcastHub {
	return &BroadcastHub{
		incoming:    incoming,
		outgoing:    outgoing,
		register:    make(chan *Connection, 1024),
		unregister:  make(chan *Connection, 1024),
		connections: make(map[*Connection]bool),
	}
}

func (bh *BroadcastHub) Run() {
	go func() {
		for {
			select {
			case c := <-bh.register:
				bh.connections[c] = true
			case c := <-bh.unregister:
				if _, ok := bh.connections[c]; ok {
					delete(bh.connections, c)
					close(c.send)
				}
			case m := <-bh.outgoing:
				for c := range bh.connections {
					select {
					case c.send <- m:
					default:
						close(c.send)
						delete(bh.connections, c)
					}
				}
			}
		}
	}()
}

func (bh *BroadcastHub) Count() int {
	return len(bh.connections)
}
