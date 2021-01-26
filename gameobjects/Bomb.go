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
	"time"

	tb "github.com/nsf/termbox-go"

	z "../common"
)

type Bomb struct {
	*Creature

	colors []tb.Attribute
}

func NewBomb(broadcast chan *z.Message, rooms z.IRooms) *Bomb {
	colors := []tb.Attribute{
		z.BoldColorRed,
		z.BoldColorGreen,
		z.BoldColorYellow,
		z.BoldColorBlue,
		z.BoldColorMagenta,
		z.BoldColorCyan,
		z.BoldColorWhite}

	return &Bomb{
		Creature: &Creature{
			GameObject: &GameObject{
				Class:     "Bomb",
				Name:      "Bomb",
				Symbol:    '▲',
				Color:     z.BoldColorWhite,
				ID:        z.UUID(),
				broadcast: broadcast,
				Paused:    true,
			},
			rooms:    rooms,
			Health:   0,
			Strength: 25,
		},
		colors: colors,
	}
}

func (b *Bomb) Run(broadcast bool) {
	if broadcast {
		msg := b.Event("Run")
		b.broadcast <- msg
	}

	go b.Loop(200, 200, b.Body)
}

func (b *Bomb) Body() {
	if b.Dead() {
		b.Explode(true)
		b.Stop(true)
		b.Delete(true)

		return
	}

	b.Battle(b.fight)
}

func (b *Bomb) Explode(broadcast bool) {
	if broadcast {
		msg := b.Event("Explode")
		b.broadcast <- msg
	}

	msg := b.Event("Sfx")
	msg.Params["Effect"] = "explode"
	b.broadcast <- msg

	for _, c := range b.colors {
		b.SetColor(c)

		time.Sleep(200 * time.Millisecond)
	}
}

func (b *Bomb) fight(opponent z.ICreature) {
	fighterName := b.GetName()
	opponentName := opponent.GetName()
	opponentSymbol := opponent.GetSymbol()
	player := opponentSymbol == '☻' || opponentSymbol == '☺'

	if player {
		hit := -1 * b.GetStrength()
		opponent.ChangeHealth(true, hit)
		b.ChangeHealth(true, hit)

		status := fighterName + " hit " + opponentName + "!"
		color := z.ColorWhite

		msg := b.Event("Announce")
		msg.Params["Status"] = status
		msg.Params["Color"] = fmt.Sprintf("%d", color)
		b.broadcast <- msg

		if opponent.Dead() {
			return
		}
	}

	opponent.Release(true)
}
