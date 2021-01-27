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
	z "./common"
	zgo "./gameobjects"
)

func (g *Game) newClient() string {
	defer g.recover()

	g.announce(true, "Pausing game", z.BoldColorWhite)

	g.pause(true, true)

	json := g.getCurrentState()

	g.announce(true, "New client", z.BoldColorWhite)

	return json
}

func (g *Game) newPlayer(broadcast bool, name, id string, opponent bool) {
	defer g.recover()

	symbol := '☻'
	color := z.BoldColorYellow

	if opponent {
		symbol = '☺'
		color = z.BoldColorMagenta
	}

	p := zgo.NewPlayer(broadcast, g.broadcast, g.id, g.config.WorldWidth, g.config.WorldHeight, g.rooms, name, id, symbol, color)

	g.players.Set(id, p)

	g.sfx(broadcast, "teleport")
}

func (g *Game) newMissle(broadcast bool, class, id string, x, y, nextX, nextY int) {
	defer g.recover()

	m := zgo.NewMissle(broadcast, g.broadcast, g.id, g.config.WorldWidth, g.config.WorldHeight, g.rooms, class, id, x, y, nextX, nextY)

	g.missles.Set(id, m)

	g.sfx(broadcast, "fire")
}
