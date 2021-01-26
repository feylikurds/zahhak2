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

const (
	HOST_IP       = "127.0.0.1"
	PORT_NUM      = "1947"
	NAME          = "Player"
	DIFFICULTY    = 30
	WORLD_WIDTH   = 20
	WORLD_HEIGHT  = 20
	NUM_MONSTERS  = 5
	NUM_HEALTHS   = 10
	NUM_STRENGTHS = 10
	NUM_TREASURES = 10
	NUM_BOMBS     = 10
	NUM_PORTALS   = 10
	CAPACITY      = 2
	VOLUME        = 10
	DYNAMIC       = true
)

const (
	MENU_WIDTH       = 30
	MENU_HEIGHT      = 5
	MAX_MSGS_DISPLAY = 9
	STATUS_LEN       = 29
	STRENGTH_LOST    = -1
)

const (
	AttrBold tb.Attribute = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
)

const (
	AttrColorDefault tb.Attribute = iota
	AttrColorBlack
	AttrColorRed
	AttrColorGreen
	AttrColorYellow
	AttrColorBlue
	AttrColorMagenta
	AttrColorCyan
	AttrColorWhite
)

const (
	ColorBlack   = AttrColorBlack & ^AttrBold
	ColorRed     = AttrColorRed & ^AttrBold
	ColorGreen   = AttrColorGreen & ^AttrBold
	ColorYellow  = AttrColorYellow & ^AttrBold
	ColorBlue    = AttrColorBlue & ^AttrBold
	ColorMagenta = AttrColorMagenta & ^AttrBold
	ColorCyan    = AttrColorCyan & ^AttrBold
	ColorWhite   = AttrColorWhite & ^AttrBold
)

const (
	BoldColorBlack   = AttrColorBlack | AttrBold
	BoldColorRed     = AttrColorRed | AttrBold
	BoldColorGreen   = AttrColorGreen | AttrBold
	BoldColorYellow  = AttrColorYellow | AttrBold
	BoldColorBlue    = AttrColorBlue | AttrBold
	BoldColorMagenta = AttrColorMagenta | AttrBold
	BoldColorCyan    = AttrColorCyan | AttrBold
	BoldColorWhite   = AttrColorWhite | AttrBold
)
