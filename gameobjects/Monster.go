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

	z "../common"
)

type Monster struct {
	*Creature

	players    z.IGameObjectMap
	player     z.ICreature
	PlayerID   string
	Difficulty int
}

func NewMonster(broadcast chan *z.Message, worldWidth, worldHeight int, rooms z.IRooms, players z.IGameObjectMap, difficulty int) *Monster {
	return &Monster{
		Creature: &Creature{
			GameObject: &GameObject{
				Class:     "Monster",
				Name:      "Monster",
				Symbol:    '☼',
				Color:     z.BoldColorRed,
				ID:        z.UUID(),
				broadcast: broadcast,
				Paused:    true,
			},
			WorldWidth:  worldWidth,
			WorldHeight: worldHeight,
			rooms:       rooms,
			Health:      100,
			Strength:    10,
		},
		players:    players,
		Difficulty: difficulty,
	}
}

func (m *Monster) Run(broadcast bool) {
	if broadcast {
		msg := m.Event("Run")
		m.broadcast <- msg
	}

	go m.Blink()
	go m.selectPlayer()

	go m.Loop(150, 200, m.Body)
}

func (m *Monster) Body() {
	if m.Dead() {
		m.Stop(true)
		m.Delete(true)

		status := m.GetName() + " was killed!"
		color := z.BoldColorMagenta

		msg := m.Event("Announce")
		msg.Params["Status"] = status
		msg.Params["Color"] = fmt.Sprintf("%d", color)
		m.broadcast <- msg

		msg = m.Event("Sfx")
		msg.Params["Effect"] = "die"
		m.broadcast <- msg

		return
	}

	if !m.Halted() {
		nextX, nextY := m.hunt()
		x, y := m.GetPosition()

		m.Move(true, x+nextX, y+nextY)
	}

	m.Battle(m.fight)
}

func (m *Monster) Blink() {
	for !m.Deleted() {
		if !m.Running() {
			time.Sleep(200 * time.Millisecond)

			continue
		}

		currentColor := m.GetColor()

		if currentColor == z.ColorRed {
			m.SetColor(z.BoldColorRed)
		} else {
			m.SetColor(z.ColorRed)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (m *Monster) selectPlayer() {
	for !m.Deleted() {
		if !m.Running() {
			time.Sleep(200 * time.Millisecond)

			continue
		}

		p, e := m.players.GetRandomValue()

		if e == nil {
			m.SetPlayer(true, p.GetID())
		} else {
			m.SetPlayer(true, "")
		}

		time.Sleep(10000 * time.Millisecond)
	}
}

func (m *Monster) SetPlayer(broadcast bool, id string) {
	m.Lock()
	defer m.Unlock()

	if broadcast {
		msg := m.Event("SetPlayer")
		msg.Params["Player"] = id
		m.broadcast <- msg
	}

	m.PlayerID = id

	if id == "" {
		m.player = nil
	} else {
		igo, _ := m.players.Get(id)
		m.player = igo.(z.ICreature)
	}
}

func (m *Monster) getPlayer() z.ICreature {
	m.RLock()
	defer m.RUnlock()

	return m.player
}

func (m *Monster) hunt() (int, int) {
	if m.player == nil {
		return z.MoveRandomly()
	} else if z.Biased(m.Difficulty) {
		return z.MoveRandomly()
	} else {
		return z.MoveToGameObject(m, m.player)
	}
}

func (m *Monster) fight(opponent z.ICreature) {
	fighterName := m.GetName()
	opponentName := opponent.GetName()
	opponentSymbol := opponent.GetSymbol()
	nonPlayer := opponentSymbol != '▲' && opponentSymbol != '☼' && opponentSymbol != '*' && opponentSymbol != '◘'

	if nonPlayer {
		hit := -1 * m.GetStrength()
		opponent.ChangeHealth(true, hit)

		status := fighterName + " attacked " + opponentName + "!"
		color := z.ColorRed

		if opponentName == "Player" {
			color = z.BoldColorRed
		}

		msg := m.Event("Announce")
		msg.Params["Status"] = status
		msg.Params["Color"] = fmt.Sprintf("%d", color)
		m.broadcast <- msg

		if opponent.Dead() {
			return
		}
	}

	opponent.Release(true)
}

func (m *Monster) LoadPlayers(players z.IGameObjectMap) {
	m.Lock()
	defer m.Unlock()

	m.players = players
}

func (m *Monster) LoadPlayer() {
	m.Lock()
	defer m.Unlock()

	id := m.PlayerID

	if id == "" {
		m.player = nil
	} else {
		igo, _ := m.players.Get(id)
		m.player = igo.(z.ICreature)
	}
}
