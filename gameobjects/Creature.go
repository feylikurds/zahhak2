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

package gameobjects

import (
	"fmt"

	z "../common"
)

type Creature struct {
	*GameObject

	WorldWidth  int
	WorldHeight int
	rooms       z.IRooms

	NextX    int
	NextY    int
	Health   int
	Strength int
	Treasure int
	Stuck    bool
}

func NewCreature(broadcast chan *z.Message, worldWidth, worldHeight int, rooms z.IRooms) *Creature {
	return &Creature{
		GameObject: &GameObject{
			Class:     "Creature",
			Name:      "Creature",
			Symbol:    'C',
			broadcast: broadcast,
			Paused:    true,
		},
		WorldWidth:  worldWidth,
		WorldHeight: worldHeight,
		rooms:       rooms,
	}
}

func (c *Creature) Next(broadcast bool, x, y int) {
	c.Lock()
	defer c.Unlock()

	if x == 0 && y == 0 {
		return
	}

	if broadcast {
		msg := c.Event("Next")
		msg.Params["X"] = fmt.Sprintf("%d", x)
		msg.Params["Y"] = fmt.Sprintf("%d", y)
		c.broadcast <- msg
	}

	c.NextX, c.NextY = x, y
}

func (c *Creature) GetNext() (int, int) {
	c.RLock()
	defer c.RUnlock()

	return c.NextX, c.NextY
}

func (c *Creature) ChangeHealth(broadcast bool, health int) {
	c.Lock()
	defer c.Unlock()

	if broadcast {
		msg := c.Event("ChangeHealth")
		msg.Params["Points"] = fmt.Sprintf("%d", health)
		c.broadcast <- msg
	}

	c.Health += health
}

func (c *Creature) GetHealth() int {
	c.RLock()
	defer c.RUnlock()

	return c.Health
}

func (c *Creature) ChangeStrength(broadcast bool, strength int) {
	c.Lock()
	defer c.Unlock()

	if broadcast {
		msg := c.Event("ChangeStrength")
		msg.Params["Points"] = fmt.Sprintf("%d", strength)
		c.broadcast <- msg
	}

	c.Strength += strength
}

func (c *Creature) GetStrength() int {
	c.RLock()
	defer c.RUnlock()

	return c.Strength
}

func (c *Creature) ChangeTreasure(broadcast bool, treasure int) {
	c.Lock()
	defer c.Unlock()

	if broadcast {
		msg := c.Event("ChangeTreasure")
		msg.Params["Points"] = fmt.Sprintf("%d", treasure)
		c.broadcast <- msg
	}

	c.Treasure += treasure
}

func (c *Creature) GetTreasure() int {
	c.RLock()
	defer c.RUnlock()

	return c.Treasure
}

func (c *Creature) Dead() bool {
	c.RLock()
	defer c.RUnlock()

	return c.Health < 0
}

func (c *Creature) Stay(broadcast bool) {
	c.Lock()
	defer c.Unlock()

	if broadcast {
		msg := c.Event("Stay")
		c.broadcast <- msg
	}

	c.Stuck = true
}

func (c *Creature) Release(broadcast bool) {
	c.Lock()
	defer c.Unlock()

	if broadcast {
		msg := c.Event("Release")
		c.broadcast <- msg
	}

	c.Stuck = false
}

func (c *Creature) Halted() bool {
	c.RLock()
	defer c.RUnlock()

	return c.Stuck
}

func (c *Creature) Move(broadcast bool, nextX, nextY int) {
	if nextX < 0 {
		nextX = c.WorldWidth - 1
	} else if nextX >= c.WorldWidth {
		nextX = 0
	}

	if nextY < 0 {
		nextY = c.WorldHeight - 1
	} else if nextY >= c.WorldHeight {
		nextY = 0
	}

	if c.rooms.HasRoom(nextX, nextY) {
		x, y := c.GetPosition()

		c.rooms.Leave(broadcast, x, y, c)
		c.rooms.Enter(broadcast, nextX, nextY, c)
		c.SetPosition(broadcast, nextX, nextY)
	}
}

func (c *Creature) Battle(fight func(z.ICreature)) {
	if c.GetStrength() > 0 {
		x, y := c.GetPosition()

		cs := c.rooms.GetCreatures(x, y)

		for _, opponent := range cs {
			if c.Dead() {
				break
			}

			oID := opponent.GetID()
			cID := c.GetID()

			if cID == oID {
				continue
			}

			opponent.Stay(true)
			fight(opponent)
		}
	}
}

func (c *Creature) LoadRooms(rooms z.IRooms) {
	c.Lock()
	defer c.Unlock()

	c.rooms = rooms
}
