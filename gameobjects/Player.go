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
	"fmt"

	tb "github.com/nsf/termbox-go"

	z "../common"
)

type Player struct {
	*Creature
}

func NewPlayer(b bool, broadcast chan *z.Message, gameID string, worldWidth, worldHeight int, rooms z.IRooms, name, id string, symbol rune, color tb.Attribute) *Player {
	if b {
		msg := z.NewMessage("Game", gameID, "NewPlayer")
		msg.Params["Name"] = name
		msg.Params["ID"] = id
		broadcast <- msg
	}

	return &Player{
		Creature: &Creature{
			GameObject: &GameObject{
				Class:     "Player",
				Name:      name,
				Symbol:    symbol,
				Color:     color,
				ID:        id,
				broadcast: broadcast,
				Paused:    true,
			},
			WorldWidth:  worldWidth,
			WorldHeight: worldHeight,
			rooms:       rooms,
			Health:      100,
			Strength:    20,
		},
	}
}

func (p *Player) Run(broadcast bool) {
	if broadcast {
		msg := p.Event("Run")
		p.broadcast <- msg
	}

	go p.Loop(125, 200, p.Body)
}

func (p *Player) Body() {
	if p.Dead() {
		p.Stop(true)
		p.Delete(true)

		status := p.GetName() + " was killed!"
		color := z.BoldColorMagenta

		msg := p.Event("Announce")
		msg.Params["Status"] = status
		msg.Params["Color"] = fmt.Sprintf("%d", color)
		p.broadcast <- msg

		msg = p.Event("Sfx")
		msg.Params["Effect"] = "die"
		p.broadcast <- msg

		return
	}

	if !p.Halted() {
		nextX, nextY := p.GetNext()
		x, y := p.GetPosition()

		p.Move(true, x+nextX, y+nextY)

		p.Collect(true)
	}

	p.Battle(p.fight)
}

func (p *Player) Collect(broadcast bool) {
	x, y := p.GetPosition()

	healths := p.rooms.GetHealths(x, y)

	for _, health := range healths {
		p.ChangeHealth(true, health.GetPoints())
		health.Delete(true)
	}

	if len(healths) > 0 {
		msg := p.Event("Announce")
		msg.Params["Status"] = "Player got health"
		msg.Params["Color"] = fmt.Sprintf("%d", z.BoldColorGreen)
		p.broadcast <- msg

		msg = p.Event("Sfx")
		msg.Params["Effect"] = "item"
		p.broadcast <- msg
	}

	strengths := p.rooms.GetStrengths(x, y)

	for _, strength := range strengths {
		p.ChangeStrength(true, strength.GetPoints())
		strength.Delete(true)
	}

	if len(strengths) > 0 {
		msg := p.Event("Announce")
		msg.Params["Status"] = "Player got strength"
		msg.Params["Color"] = fmt.Sprintf("%d", z.BoldColorCyan)
		p.broadcast <- msg

		msg = p.Event("Sfx")
		msg.Params["Effect"] = "item"
		p.broadcast <- msg
	}

	treasures := p.rooms.GetTreasures(x, y)

	for _, treasure := range treasures {
		p.ChangeTreasure(true, treasure.GetPoints())
		treasure.Delete(true)
	}

	if len(treasures) > 0 {
		msg := p.Event("Announce")
		msg.Params["Status"] = "Player got treasure"
		msg.Params["Color"] = fmt.Sprintf("%d", z.BoldColorBlue)
		p.broadcast <- msg

		msg = p.Event("Sfx")
		msg.Params["Effect"] = "item"
		p.broadcast <- msg
	}
}

func (p *Player) fight(opponent z.ICreature) {
	fighterName := p.GetName()
	opponentName := opponent.GetName()
	opponentSymbol := opponent.GetSymbol()
	nonObject := opponentSymbol != '▲' && opponentSymbol != '◘' && opponentSymbol != '*'

	if nonObject {
		hit := -1 * p.GetStrength()
		opponent.ChangeHealth(true, hit)
		p.ChangeStrength(true, z.STRENGTH_LOST)

		status := fighterName + " attacked " + opponentName + "!"
		color := z.BoldColorYellow

		msg := p.Event("Announce")
		msg.Params["Status"] = status
		msg.Params["Color"] = fmt.Sprintf("%d", color)
		p.broadcast <- msg

		if opponent.Dead() {
			return
		}
	}

	opponent.Release(true)
}
