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

package gameobjects

import (
	"strconv"
	"sync"

	z "../common"
)

type Room struct {
	*GameObject

	GameObjects map[z.IGameObject]z.IGameObject
	Capacity    int

	sync.RWMutex
}

func NewRoom(broadcast chan *z.Message, capacity, x, y int) *Room {
	gom := make(map[z.IGameObject]z.IGameObject, capacity)

	return &Room{
		GameObject: &GameObject{
			Class:     "Room",
			Name:      "Room",
			Symbol:    'R',
			Color:     z.ColorBlack,
			ID:        z.UUID(),
			X:         x,
			Y:         y,
			broadcast: broadcast,
			Paused:    true,
		},
		GameObjects: gom,
		Capacity:    capacity,
	}
}

func (r *Room) HasRoomForTwo() bool {
	r.RLock()
	defer r.RUnlock()

	n := r.Capacity - len(r.GameObjects)

	if n >= 2 {
		return true
	}

	return false
}

func (r *Room) HasRoom() bool {
	r.RLock()
	defer r.RUnlock()

	n := len(r.GameObjects)

	if n+1 > r.Capacity {
		return false
	}

	return true
}

func (r *Room) Enter(broadcast bool, g z.IGameObject) {
	r.Lock()
	defer r.Unlock()

	if broadcast {
		msg := r.Event("Enter")
		msg.Params["Class"] = g.GetClass()
		msg.Params["ID"] = g.GetID()
		msg.Params["X"] = strconv.Itoa(r.X)
		msg.Params["Y"] = strconv.Itoa(r.Y)
		r.broadcast <- msg
	}

	r.GameObjects[g] = g
}

func (r *Room) Leave(broadcast bool, g z.IGameObject) {
	r.Lock()
	defer r.Unlock()

	if broadcast {
		msg := r.Event("Leave")
		msg.Params["Class"] = g.GetClass()
		msg.Params["ID"] = g.GetID()
		msg.Params["X"] = strconv.Itoa(r.X)
		msg.Params["Y"] = strconv.Itoa(r.Y)
		r.broadcast <- msg
	}

	delete(r.GameObjects, g)
}

func (r *Room) GetGameObjects() []z.IGameObject {
	r.Lock()
	defer r.Unlock()

	gos := []z.IGameObject{}

	for g := range r.GameObjects {
		if g.Deleted() {
			delete(r.GameObjects, g)

			continue
		}

		gos = append(gos, g)
	}

	return gos
}

func (r *Room) GetCreatures() []z.ICreature {
	r.RLock()
	defer r.RUnlock()

	cs := []z.ICreature{}

	for g := range r.GameObjects {
		c, ok := g.(z.ICreature)

		if ok {
			cs = append(cs, c)
		}
	}

	return cs
}

func (r *Room) GetHealths() []z.IHealth {
	r.RLock()
	defer r.RUnlock()

	healths := []z.IHealth{}

	for g := range r.GameObjects {
		_, ok := g.(*Health)

		if ok {
			healths = append(healths, g.(z.IHealth))
		}
	}

	return healths
}

func (r *Room) GetStrengths() []z.IStrength {
	r.RLock()
	defer r.RUnlock()

	strengths := []z.IStrength{}

	for g := range r.GameObjects {
		_, ok := g.(*Strength)

		if ok {
			strengths = append(strengths, g.(z.IStrength))
		}
	}

	return strengths
}

func (r *Room) GetTreasures() []z.ITreasure {
	r.RLock()
	defer r.RUnlock()

	treasures := []z.ITreasure{}

	for g := range r.GameObjects {
		_, ok := g.(*Treasure)

		if ok {
			treasures = append(treasures, g.(z.ITreasure))
		}
	}

	return treasures
}
