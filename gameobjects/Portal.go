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

type Portal struct {
	*Creature
}

func NewPortal(broadcast chan *z.Message, worldWidth, worldHeight int, rooms z.IRooms) *Portal {
	return &Portal{
		Creature: &Creature{
			GameObject: &GameObject{
				Class:     "Portal",
				Name:      "Portal",
				Symbol:    '◘',
				Color:     z.BoldColorYellow,
				ID:        z.UUID(),
				broadcast: broadcast,
				Paused:    true,
			},
			WorldWidth:  worldWidth,
			WorldHeight: worldHeight,
			rooms:       rooms,
			Health:      1000,
			Strength:    1000,
		},
	}
}

func (p *Portal) Run(broadcast bool) {
	if broadcast {
		msg := p.Event("Run")
		p.broadcast <- msg
	}

	go p.Loop(100, 200, p.Body)
}

func (p *Portal) Body() {
	p.Battle(p.fight)
}

func (p *Portal) fight(opponent z.ICreature) {
	opponentName := opponent.GetName()
	opponentSymbol := opponent.GetSymbol()
	player := opponentSymbol == '☻' || opponentSymbol == '☺'

	if player {
		status := opponentName + " teleported!"
		color := z.BoldColorYellow

		msg := p.Event("Announce")
		msg.Params["Status"] = status
		msg.Params["Color"] = fmt.Sprintf("%d", color)
		p.broadcast <- msg

		msg = p.Event("Sfx")
		msg.Params["Effect"] = "teleport"
		p.broadcast <- msg

		p.Teleport(opponent)
	}

	opponent.Release(true)
}

func (p *Portal) Teleport(c z.ICreature) {
	for {
		x := z.RandomNumber(0, p.WorldWidth-1)
		y := z.RandomNumber(0, p.WorldHeight-1)

		if p.rooms.HasRoomForTwo(x, y) {
			c.Move(true, x, y)

			break
		}
	}
}
