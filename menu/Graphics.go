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

package menu

import (
	"errors"

	rw "github.com/mattn/go-runewidth"
	tb "github.com/nsf/termbox-go"

	z "../common"
)

type Graphics struct {
	Ratio        float64
	ScreenWidth  int
	ScreenHeight int
	UnitX        int
	UnitY        int
	CenterX      int
	CenterY      int
}

func NewGraphics(ratio float64) *Graphics {
	screenWidth, screenHeight := tb.Size()
	unitX := int(float64(screenWidth) / ratio)
	unitY := int(float64(screenHeight) / ratio)
	centerX := int(float64(screenWidth) / 2)
	centerY := int(float64(screenHeight) / 2)

	return &Graphics{
		Ratio:        ratio,
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		UnitX:        unitX,
		UnitY:        unitY,
		CenterX:      centerX,
		CenterY:      centerY,
	}
}

func (g *Graphics) SetCursor(x, y int) {
	tb.SetCursor(x, y)
}

func (g *Graphics) Flush() {
	tb.Flush()
}

func (g *Graphics) Resize(ratio float64) {
	g.Ratio = ratio
	g.resize()
}

func (g *Graphics) resize() {
	g.ScreenWidth, g.ScreenHeight = tb.Size()
	g.UnitX = int(float64(g.ScreenWidth) / g.Ratio)
	g.UnitY = int(float64(g.ScreenHeight) / g.Ratio)
	g.CenterX = int(float64(g.ScreenWidth) / 2)
	g.CenterY = int(float64(g.ScreenHeight) / 2)
}

func (g *Graphics) GetUnits() (int, int) {
	return g.UnitX, g.UnitY
}

func (g *Graphics) GetCenter() (int, int) {
	return g.CenterX, g.CenterY
}

func (g *Graphics) BlankScreen() {
	tb.Clear(tb.ColorBlack, tb.ColorBlack)
}

func (g *Graphics) tbprint(x, y int, fg, bg tb.Attribute, msg string) {
	for _, c := range msg {
		tb.SetCell(x, y, c, fg, bg)
		x += rw.RuneWidth(c)
	}
}

func (g *Graphics) Print(x, y int, msg string, color tb.Attribute) {
	g.tbprint(x, y, color, tb.ColorDefault, msg)
}

func (g *Graphics) PrintFgBg(x, y int, msg string, fgColor tb.Attribute, bgColor tb.Attribute) {
	g.tbprint(x, y, fgColor, bgColor, msg)
}

func (g *Graphics) DrawUnit(xU int, yU int, symbol string, color tb.Attribute) {
	g.DrawUnitFgBg(xU, yU, symbol, color, z.ColorBlack)
}

func (g *Graphics) DrawUnitFgBg(xU int, yU int, symbol string, fgColor tb.Attribute, bgColor tb.Attribute) {
	for y := yU; y < yU+g.UnitY; y++ {
		for x := xU; x < xU+g.UnitX; x++ {
			g.PrintFgBg(x, y, symbol, fgColor, bgColor)
		}
	}
}

func (g *Graphics) DrawBackground(symbol string, fgColor tb.Attribute, ratio float64) {
	g.DrawBackgroundFgBg(symbol, fgColor, z.ColorBlack, ratio)
}

func (g *Graphics) DrawBackgroundFgBg(symbol string, fgColor tb.Attribute, bgColor tb.Attribute, ratio float64) {
	oldRatio := g.Ratio
	g.Ratio = ratio
	g.resize()

	cX, cY := g.GetCenter()
	uX, uY := g.GetUnits()

	for x := cX - uX; x < cX; x++ {
		for y := cY - uY; y < cY; y++ {
			g.DrawUnitFgBg(x, y, symbol, fgColor, bgColor)
		}
	}
	g.Ratio = oldRatio
	g.resize()
}

func (g *Graphics) ReadChar() (rune, error) {
	for {
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			switch ev.Key {
			case tb.KeyEnter:
				return '\n', nil
			case tb.KeySpace:
				return ' ', nil
			default:
				return ev.Ch, nil
			}
		}
	}

	return '\000', errors.New("No key pressed")
}

func (g *Graphics) Readline(x, y int) string {
	var line []rune

	g.SetCursor(x, y)
	g.Flush()

	for {
		c, e := g.ReadChar()

		if e == nil {
			if c == '\n' {
				break
			}

			if c > 31 && c < 127 {
				line = append(line, c)
			}

			g.Print(x, y, string(c), z.BoldColorWhite)
			x++
			g.SetCursor(x, y)
			g.Flush()
		}
	}

	s := string(line)

	return s
}
