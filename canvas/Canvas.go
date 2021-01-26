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

package canvas

import (
	"fmt"
	"strconv"
	"time"

	rw "github.com/mattn/go-runewidth"
	tb "github.com/nsf/termbox-go"

	z "../common"
)

type Canvas struct {
	rooms    z.IRooms
	statuses z.IRing

	menuWidth      int
	numMsgsDisplay int
	worldWidth     int
	worldHeight    int
	screenWidth    int
	screenHeight   int
}

func NewCanvas(c *z.Config, rooms z.IRooms, statuses z.IRing) *Canvas {
	defer tb.Flush()
	tb.Clear(tb.ColorBlack, tb.ColorBlack)

	screenWidth := c.WorldWidth*c.Capacity + c.MenuWidth
	screenHeight := c.WorldHeight + c.MenuHeight

	return &Canvas{
		rooms:          rooms,
		statuses:       statuses,
		menuWidth:      c.MenuWidth,
		numMsgsDisplay: c.NumMsgsDisplay,
		worldWidth:     c.WorldWidth,
		worldHeight:    c.WorldHeight,
		screenWidth:    screenWidth,
		screenHeight:   screenHeight}
}

func (c *Canvas) Draw(numHealths, numStrengths, numTreasures, totalTreasures int) {
	tb.Sync()
	defer tb.Flush()

	c.paint()

	c.stats(numHealths, numStrengths, numTreasures, totalTreasures)

	c.overlay()
}

func (c *Canvas) print(x, y int, msg string, color tb.Attribute) {
	c.tbprint(x, y, color, z.ColorBlack, msg)
}

func (c *Canvas) tbprint(x, y int, fg, bg tb.Attribute, msg string) {
	for _, c := range msg {
		tb.SetCell(x, y, c, fg, bg)
		x += rw.RuneWidth(c)
	}
}

func (c *Canvas) stats(numHealths, numStrengths, numTreasures, totalTreasures int) {
	now := time.Now()
	t := fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
	col := c.worldWidth + c.menuWidth - len(t)
	row := c.screenHeight - 1
	c.print(col, row, t, z.BoldColorWhite)

	row = 2
	col = c.worldWidth + 2
	s := c.entryTextValue("Health", strconv.Itoa(numHealths))
	c.print(col, row, s, z.BoldColorGreen)

	s = c.entryTextValue("Strength", strconv.Itoa(numStrengths))
	c.print(col, row+1, s, z.BoldColorCyan)

	s = c.entryTextValue("Treasure", strconv.Itoa(numTreasures)+"/"+strconv.Itoa(totalTreasures))
	c.print(col, row+2, s, z.BoldColorBlue)

	row = 2 + 3 + c.numMsgsDisplay
	statuses := c.statuses.Values()

	for _, s := range statuses {
		status, ok := s.(*z.Status)

		if !ok {
			continue
		}

		text := c.entryTextPad(status.Text)
		c.print(col, row, text, status.Color)

		row--
	}
}

func (c *Canvas) overlay() {
	t := "Zahhak2 by Aryo Pehlewan feylikurds@gmail.com Copyright 2016 License GPLv3"
	s := c.entryTextLen(t, len(t))
	c.print(0, 0, s, z.BoldColorWhite)

	col := 0
	row := c.screenHeight - 1

	t = "☻ : Player"
	s = c.entryText(t)
	c.print(col, row-1, s, z.BoldColorYellow)
	t = "☼ : Monster"
	s = c.entryText(t)
	c.print(col, row, s, z.BoldColorRed)

	t = "H : Health"
	s = c.entryText(t)
	c.print(col+12, row-1, s, z.BoldColorGreen)
	t = "S : Strength"
	s = c.entryText(t)
	c.print(col+12, row, s, z.BoldColorCyan)

	t = "T : Treasure"
	s = c.entryText(t)
	c.print(col+25, row-1, s, z.BoldColorBlue)
	t = "☺ : Opponent"
	s = c.entryText(t)
	c.print(col+25, row, s, z.BoldColorMagenta)

	t = "▲ : Bomb"
	s = c.entryText(t)
	c.print(col+38, row-1, s, z.BoldColorWhite)
	t = "◘ : Portal"
	s = c.entryText(t)
	c.print(col+38, row, s, z.BoldColorYellow)

	t = "Type 'zahak2 help' for options"
	col = c.worldWidth + c.menuWidth - len(t)
	row = c.screenHeight - 2
	c.print(col, row, t, z.ColorWhite)

	col = c.worldWidth + 2
	row = c.worldHeight + 1

	s = c.entryText("Enter: Pause")
	c.print(col, row-3, s, z.BoldColorYellow)
	s = c.entryText("Esc/Q: Quit")
	c.print(col, row-2, s, z.BoldColorYellow)
	s = c.entryText("Arrows/Left mouse: Move")
	c.print(col, row-1, s, z.BoldColorYellow)
	s = c.entryText("Space/Right mouse: Shoot")
	c.print(col, row, s, z.BoldColorYellow)

	row = 1

	for col := 0; col < c.screenWidth; col++ {
		c.print(col, row, "═", z.BoldColorYellow)
	}

	row = c.worldHeight + 2

	for col := 0; col < c.screenWidth; col++ {
		c.print(col, row, "═", z.BoldColorYellow)
	}

	row = 5

	for col := c.worldWidth + 1; col < c.screenWidth; col++ {
		c.print(col, row, "═", z.BoldColorYellow)
	}

	row = 2 + 4 + c.numMsgsDisplay

	for col := c.worldWidth + 1; col < c.screenWidth; col++ {
		c.print(col, row, "═", z.BoldColorYellow)
	}

	col = 0

	for row := 1; row < c.screenHeight-2; row++ {
		c.print(col, row, "║", z.BoldColorYellow)
	}

	col = c.worldWidth + 1

	for row := 1; row < c.screenHeight-2; row++ {
		c.print(col, row, "║", z.BoldColorYellow)
	}

	row = 1
	col = 0
	c.print(col, row, "╔", z.BoldColorYellow)
	row = c.worldHeight + 2
	col = 0
	c.print(col, row, "╚", z.BoldColorYellow)
	row = 1
	col = c.worldWidth + 1
	c.print(col, row, "╦", z.BoldColorYellow)
	row = c.worldHeight + 2
	col = c.worldWidth + 1
	c.print(col, row, "╩", z.BoldColorYellow)
	row = 5
	col = c.worldWidth + 1
	c.print(col, row, "╠", z.BoldColorYellow)
	row = c.worldHeight - 3
	col = c.worldWidth + 1
	c.print(col, row, "╠", z.BoldColorYellow)
}

func (c *Canvas) paint() {
	for y := 0; y < c.worldHeight; y++ {
		for x := 0; x < c.worldWidth; x++ {
			gos := c.rooms.GetGameObjects(x, y)
			l := len(gos)

			if l == 0 {
				tb.SetCell(x+1, y+2, '.', z.BoldColorWhite, z.ColorBlack)
			} else {
				for i, g := range gos {
					tb.SetCell(x+1+i, y+2, g.GetSymbol(), g.GetColor(), z.ColorBlack)
				}
			}
		}
	}
}

func (c *Canvas) entryText(text string) string {
	return c.status(text, z.STATUS_LEN)
}

func (c *Canvas) entryTextPad(text string) string {
	tp := c.padRight(text, " ", z.STATUS_LEN)

	return c.status(tp, z.STATUS_LEN)
}

func (c *Canvas) padRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func (c *Canvas) entryTextValue(text, value string) string {
	tv := fmt.Sprintf("%-8s", text) + " = " + fmt.Sprintf("%-6s", value)

	return c.status(tv, z.STATUS_LEN)
}

func (c *Canvas) entryTextLen(text string, length int) string {
	return c.status(text, length)
}

func (c *Canvas) status(text string, length int) string {
	symbol := text

	if len(text) > length {
		symbol = text[0:length]
	}

	return symbol
}
