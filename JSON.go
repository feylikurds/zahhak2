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

package main

import (
	"encoding/json"
	"os"

	tb "github.com/nsf/termbox-go"

	z "./common"
	zgo "./gameobjects"
)

func (g *Game) jsonGOM(session, class string, gom *GameObjectMap, mp []string) {

	for _, o := range mp {
		igo, e := g.jsonGO(class, o)

		if e != nil {
			tb.Close()

			println("Error: could not load " + class)

			os.Exit(1)
		}

		gom.Set(igo.GetID(), igo)
		x, y := igo.GetPosition()
		g.rooms.Enter(false, x, y, igo)
	}
}

func (g *Game) jsonGO(class string, o string) (z.IGameObject, error) {
	var igo z.IGameObject

	switch class {
	case "Players":
		player := &zgo.Player{}

		if e := json.Unmarshal([]byte(o[:]), player); e != nil {
			return nil, e
		}

		igo = player

	case "Monsters":
		monster := &zgo.Monster{}

		if e := json.Unmarshal([]byte(o[:]), monster); e != nil {
			return nil, e
		}

		igo = monster

	case "Bombs":
		bomb := &zgo.Bomb{}

		if e := json.Unmarshal([]byte(o[:]), bomb); e != nil {
			return nil, e
		}

		igo = bomb

	case "Portals":
		portal := &zgo.Portal{}

		if e := json.Unmarshal([]byte(o[:]), portal); e != nil {
			return nil, e
		}

		igo = portal

	case "Missles":
		missle := &zgo.Missle{}

		if e := json.Unmarshal([]byte(o[:]), missle); e != nil {
			return nil, e
		}

		igo = missle

	case "Healths":
		health := &zgo.Health{}

		if e := json.Unmarshal([]byte(o[:]), health); e != nil {
			return nil, e
		}

		igo = health

	case "Strengths":
		strength := &zgo.Strength{}

		if e := json.Unmarshal([]byte(o[:]), strength); e != nil {
			return nil, e
		}

		igo = strength

	case "Treasures":
		treasure := &zgo.Treasure{}

		if e := json.Unmarshal([]byte(o[:]), treasure); e != nil {
			return nil, e
		}

		igo = treasure
	}

	return igo, nil
}

func (g *Game) getCurrentState() string {
	m := g.Event("CurrentState")

	bs, _ := json.Marshal(g.config)
	s := string(bs)
	m.Params["Config"] = s
	m.Params["Session"] = g.session

	players := g.gomJSON(g.players)
	m.MultiParams["Players"] = players
	monsters := g.gomJSON(g.monsters)
	m.MultiParams["Monsters"] = monsters

	bombs := g.gomJSON(g.bombs)
	m.MultiParams["Bombs"] = bombs
	portals := g.gomJSON(g.portals)
	m.MultiParams["Portals"] = portals
	missles := g.gomJSON(g.missles)
	m.MultiParams["Missles"] = missles

	healths := g.gomJSON(g.healths)
	m.MultiParams["Healths"] = healths
	strengths := g.gomJSON(g.strengths)
	m.MultiParams["Strengths"] = strengths
	treasures := g.gomJSON(g.treasures)
	m.MultiParams["Treasures"] = treasures

	bs, _ = json.Marshal(m)
	s = string(bs)

	return s
}

func (g *Game) gomJSON(gom *GameObjectMap) []string {
	jsonArray := []string{}
	gos := gom.GetValues()

	for _, igo := range gos {
		bs := g.igoJSON(igo)
		s := string(bs)
		jsonArray = append(jsonArray, s)
	}

	return jsonArray
}

func (g *Game) igoJSON(igo z.IGameObject) []byte {
	bs := []byte{}

	switch igo.GetClass() {
	case "Player":
		p := igo.(*zgo.Player)
		bs, _ = json.Marshal(p)
		o := &zgo.Player{}
		json.Unmarshal(bs, o)
		o.Symbol = '☺'
		o.Color = z.BoldColorMagenta
		bs, _ = json.Marshal(o)

	case "Monster":
		o := igo.(*zgo.Monster)
		bs, _ = json.Marshal(o)

	case "Bomb":
		o := igo.(*zgo.Bomb)
		bs, _ = json.Marshal(o)

	case "Portal":
		o := igo.(*zgo.Portal)
		bs, _ = json.Marshal(o)

	case "Missle":
		o := igo.(*zgo.Missle)
		bs, _ = json.Marshal(o)

	case "Health":
		o := igo.(*zgo.Health)
		bs, _ = json.Marshal(o)

	case "Strength":
		o := igo.(*zgo.Strength)
		bs, _ = json.Marshal(o)

	case "Treasure":
		o := igo.(*zgo.Treasure)
		bs, _ = json.Marshal(o)
	}

	return bs
}
