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

package gameobjects

import (
	"fmt"
	"strconv"
	"time"

	tb "github.com/nsf/termbox-go"

	z "../common"
)

type Missle struct {
	*Creature

	NextX  int
	NextY  int
	colors []tb.Attribute
}

func NewMissle(b bool, broadcast chan *z.Message, gameID string, worldWidth, worldHeight int, rooms z.IRooms, class, id string, x, y, nextX, nextY int) *Missle {
	if b {
		msg := z.NewMessage("Game", gameID, "NewMissle")
		msg.Params["Creature"] = class
		msg.Params["ID"] = id
		msg.Params["X"] = strconv.Itoa(x)
		msg.Params["Y"] = strconv.Itoa(y)
		msg.Params["NextX"] = strconv.Itoa(nextX)
		msg.Params["NextY"] = strconv.Itoa(nextY)
		broadcast <- msg
	}

	colors := []tb.Attribute{
		z.BoldColorRed,
		z.BoldColorGreen,
		z.BoldColorYellow,
		z.BoldColorBlue,
		z.BoldColorMagenta,
		z.BoldColorCyan,
		z.BoldColorWhite}

	return &Missle{
		Creature: &Creature{
			GameObject: &GameObject{
				Class:     "Missle",
				Name:      "Missle",
				Symbol:    '*',
				Color:     z.BoldColorWhite,
				X:         x,
				Y:         y,
				ID:        id,
				broadcast: broadcast,
				Paused:    true,
			},
			WorldWidth:  worldWidth,
			WorldHeight: worldHeight,
			rooms:       rooms,
			Health:      0,
			Strength:    50,
		},
		NextX:  nextX,
		NextY:  nextY,
		colors: colors,
	}
}

func (m *Missle) Run(broadcast bool) {
	if broadcast {
		msg := m.Event("Run")
		m.broadcast <- msg
	}

	go m.Flash()

	go m.Loop(50, 200, m.Body)
}

func (m *Missle) Body() {
	if m.Dead() {
		m.Stop(true)
		m.Delete(true)

		return
	}

	nextX, nextY := m.GetNext()

	if nextX == 0 && nextY == 0 {
		m.Stop(true)
		m.Delete(true)

		return
	}

	currentX, currentY := m.GetPosition()

	m.Move(true, currentX+nextX, currentY+nextY)

	m.Battle(m.fight)
}

func (m *Missle) Move(broadcast bool, nextX, nextY int) {
	if nextX < 0 || nextX >= m.WorldWidth || nextY < 0 || nextY >= m.WorldHeight {
		m.Stop(broadcast)
		m.Delete(broadcast)

		return
	}

	if m.rooms.HasRoom(nextX, nextY) {
		x, y := m.GetPosition()
		m.rooms.Leave(broadcast, x, y, m)
		m.rooms.Enter(broadcast, nextX, nextY, m)
		m.SetPosition(broadcast, nextX, nextY)
	}
}

func (m *Missle) GetNext() (int, int) {
	m.RLock()
	defer m.RUnlock()

	return m.NextX, m.NextY
}

func (m *Missle) Flash() {
	for !m.Deleted() {
		if !m.Running() {
			time.Sleep(200 * time.Millisecond)

			continue
		}

		currentColor := m.GetColor()

		if currentColor == m.colors[len(m.colors)-1] {
			m.SetColor(m.colors[0])
		} else {
			index := 0

			for i, c := range m.colors {
				if c == currentColor {
					index = i

					break
				}
			}

			m.SetColor(m.colors[index+1])
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (m *Missle) fight(opponent z.ICreature) {
	fighterName := m.GetName()
	opponentName := opponent.GetName()

	hit := -1 * m.GetStrength()
	opponent.ChangeHealth(true, hit)
	m.ChangeHealth(true, hit)

	status := fighterName + " hit " + opponentName + "!"
	color := z.ColorWhite

	msg := m.Event("Announce")
	msg.Params["Status"] = status
	msg.Params["Color"] = fmt.Sprintf("%d", color)
	m.broadcast <- msg

	if opponent.Dead() {
		return
	}

	opponent.Release(true)
}
