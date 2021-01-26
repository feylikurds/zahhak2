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

package common

import (
	tb "github.com/nsf/termbox-go"
)

type Config struct {
	Multiplayer bool
	Server      bool
	Host        string
	Port        string
	Name        string

	Difficulty   int
	WorldWidth   int
	WorldHeight  int
	Capacity     int
	NumMonsters  int
	NumHealths   int
	NumStrengths int
	NumTreasures int
	NumBombs     int
	NumPortals   int

	Volume         int
	Dynamic        bool
	MenuWidth      int
	MenuHeight     int
	NumMsgsDisplay int
}

func NewConfig() *Config {
	return &Config{
		Host: HOST_IP,
		Port: PORT_NUM,

		Difficulty:   DIFFICULTY,
		WorldWidth:   WORLD_WIDTH,
		WorldHeight:  WORLD_HEIGHT,
		Capacity:     CAPACITY,
		NumMonsters:  NUM_MONSTERS,
		NumHealths:   NUM_HEALTHS,
		NumStrengths: NUM_STRENGTHS,
		NumTreasures: NUM_TREASURES,
		NumBombs:     NUM_BOMBS,
		NumPortals:   NUM_PORTALS,

		Volume:         VOLUME,
		Dynamic:        DYNAMIC,
		MenuWidth:      MENU_WIDTH,
		MenuHeight:     MENU_HEIGHT,
		NumMsgsDisplay: MAX_MSGS_DISPLAY}
}

func (c *Config) Init() {
	if c.Dynamic {
		c.WorldWidth, c.WorldHeight = tb.Size()
		c.WorldWidth, c.WorldHeight = c.WorldWidth-c.MenuWidth, c.WorldHeight-c.MenuHeight

		w := float64(c.WorldWidth * c.WorldHeight)
		d := float64(WORLD_WIDTH * WORLD_HEIGHT)

		if w > d {
			f := func(constant int) int {
				n := (w * (float64(constant) / d) * 0.50) + float64(constant)

				return int(n)
			}

			c.NumMonsters = f(NUM_MONSTERS)
			c.NumHealths = f(NUM_HEALTHS)
			c.NumStrengths = f(NUM_STRENGTHS)
			c.NumBombs = f(NUM_BOMBS)
			c.NumPortals = f(NUM_PORTALS)
		}
	}

	c.NumMsgsDisplay = c.WorldHeight - 4 - 5
}
