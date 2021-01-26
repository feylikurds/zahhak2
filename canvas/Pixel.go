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

package canvas

import (
	tb "github.com/nsf/termbox-go"

	z "../common"
)

type Pixel struct {
	Symbol rune
	Color  tb.Attribute
}

func NewPixel() *Pixel {
	return &Pixel{
		Symbol: ' ',
		Color:  z.ColorBlack}
}
