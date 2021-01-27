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
	"fmt"
	"os"

	tb "github.com/nsf/termbox-go"

	z "./common"
	zgo "./gameobjects"
)

func (g *Game) client() {
	name := g.config.Name
	m, ok := <-g.gameManager.networkManager.Messages

	if !ok {
		tb.Close()

		println("Error: could not connect to server")

		os.Exit(1)
	}

	p := m.Params["Config"]
	config := &z.Config{}
	bs := []byte(p[:])
	e := json.Unmarshal(bs, config)

	if e != nil {
		tb.Close()

		println("Error: could not get Config")

		os.Exit(1)
	}

	g.config = config
	g.config.Name = name
	g.session = m.Params["Session"]
	g.gameManager.SetSession(g.session)
	g.display = true

	g.createWorld()

	mp := m.MultiParams["Strengths"]
	g.jsonGOM(g.session, "Strengths", g.strengths, mp)
	s := fmt.Sprintf("Initializing %d strengths", g.strengths.Len())
	g.announce(false, s, z.BoldColorWhite)
	mp = m.MultiParams["Treasures"]
	g.jsonGOM(g.session, "Treasures", g.treasures, mp)
	s = fmt.Sprintf("Initializing %d treasures", g.treasures.Len())
	g.announce(false, s, z.BoldColorWhite)
	mp = m.MultiParams["Healths"]
	g.jsonGOM(g.session, "Healths", g.healths, mp)
	s = fmt.Sprintf("Initializing %d healths", g.healths.Len())
	g.announce(false, s, z.BoldColorWhite)

	mp = m.MultiParams["Missles"]
	g.jsonGOM(g.session, "Missles", g.missles, mp)
	g.loadRooms(g.missles)
	g.prepareMissles(g.missles)
	s = fmt.Sprintf("Initializing %d missles", g.missles.Len())
	g.announce(false, s, z.BoldColorWhite)
	mp = m.MultiParams["Portals"]
	g.jsonGOM(g.session, "Portals", g.portals, mp)
	g.loadRooms(g.portals)
	s = fmt.Sprintf("Initializing %d portals", g.portals.Len())
	g.announce(false, s, z.BoldColorWhite)
	mp = m.MultiParams["Bombs"]
	g.jsonGOM(g.session, "Bombs", g.bombs, mp)
	g.loadRooms(g.bombs)
	s = fmt.Sprintf("Initializing %d bombs", g.bombs.Len())
	g.announce(false, s, z.BoldColorWhite)

	mp = m.MultiParams["Monsters"]
	g.jsonGOM(g.session, "Monsters", g.monsters, mp)
	g.loadRooms(g.monsters)
	g.prepareMonsters(g.monsters)
	s = fmt.Sprintf("Initializing %d monsters", g.monsters.Len())
	g.announce(false, s, z.BoldColorWhite)
	mp = m.MultiParams["Players"]
	g.jsonGOM(g.session, "Players", g.players, mp)
	g.loadRooms(g.players)
	s = fmt.Sprintf("Initializing %d players", g.players.Len())
	g.announce(false, s, z.BoldColorWhite)

	g.announce(false, "Multiplayer mode", z.BoldColorWhite)
	g.announce(false, "Connecting to game", z.BoldColorWhite)
	g.announce(false, "Port: "+g.config.Port, z.BoldColorWhite)
	g.announce(false, "IP: "+g.config.Host, z.BoldColorWhite)

	id := z.UUID()
	x, y := g.randomFreePlace()

	g.newPlayer(true, g.config.Name, id, false)

	igo, _ := g.players.Get(id)
	g.player = igo.(*zgo.Player)

	g.rooms.Enter(true, x, y, igo)

	msg := g.Event("Run")
	msg.Class = "Player"
	msg.ID = id
	g.broadcast <- msg

	g.announce(false, "Initializing player", z.BoldColorWhite)

	g.announce(false, "Finished initializing", z.BoldColorWhite)

	g.pause(true, false)
}
