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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	tb "github.com/nsf/termbox-go"

	z "../common"
)

type GameObject struct {
	sync.RWMutex

	Paused  bool
	Removed bool

	Class  string
	Name   string
	Symbol rune
	Color  tb.Attribute
	X      int
	Y      int
	ID     string

	broadcast chan *z.Message
}

func NewGameObject(broadcast chan *z.Message) *GameObject {
	return &GameObject{
		Class:     "GameObject",
		Name:      "GameObject",
		Symbol:    'G',
		broadcast: broadcast,
		Paused:    true,
	}
}

func (g *GameObject) Start(broadcast bool) {
	g.Lock()
	defer g.Unlock()

	if broadcast {
		msg := g.Event("Start")
		g.broadcast <- msg
	}

	g.Paused = false
}

func (g *GameObject) Stop(broadcast bool) {
	g.Lock()
	defer g.Unlock()

	if broadcast {
		msg := g.Event("Stop")
		g.broadcast <- msg
	}

	g.Paused = true
}

func (g *GameObject) Running() bool {
	g.RLock()
	defer g.RUnlock()

	return !g.Paused
}

func (g *GameObject) Run(broadcast bool) {
	if broadcast {
		msg := g.Event("Run")
		g.broadcast <- msg
	}

	go g.Loop(200, 200, g.Body)
}

func (g *GameObject) Body() {
}

func (g *GameObject) Loop(sleepMain, sleepPaused int, body func()) {
	for !g.Deleted() {
		if !g.Running() {
			time.Sleep(time.Duration(sleepPaused) * time.Millisecond)

			continue
		}

		prev := time.Now()

		body()

		now := time.Now()
		diff := now.Sub(prev)
		wait := time.Duration(sleepMain)*time.Millisecond - diff
		time.Sleep(wait)
	}
}

func (g *GameObject) Delete(broadcast bool) {
	g.Lock()
	defer g.Unlock()

	if broadcast {
		msg := g.Event("Delete")
		g.broadcast <- msg
	}

	g.Removed = true
}

func (g *GameObject) Deleted() bool {
	g.RLock()
	defer g.RUnlock()

	return g.Removed
}

func (g *GameObject) SetID(broadcast bool, id string) {
	g.Lock()
	defer g.Unlock()

	if broadcast {
		msg := g.Event("SetID")
		msg.Params["ID"] = id
		g.broadcast <- msg
	}

	g.ID = id
}

func (g *GameObject) GetID() string {
	g.RLock()
	defer g.RUnlock()

	return g.ID
}

func (g *GameObject) SetPosition(broadcast bool, x, y int) {
	g.Lock()
	defer g.Unlock()

	if broadcast {
		msg := g.Event("SetPosition")
		msg.Params["X"] = fmt.Sprintf("%d", x)
		msg.Params["Y"] = fmt.Sprintf("%d", y)
		g.broadcast <- msg
	}

	g.X, g.Y = x, y
}

func (g *GameObject) GetPosition() (int, int) {
	g.RLock()
	defer g.RUnlock()

	return g.X, g.Y
}

func (g *GameObject) SetName(broadcast bool, name string) {
	g.Lock()
	defer g.Unlock()

	if broadcast {
		msg := g.Event("SetName")
		msg.Params["Name"] = name
		g.broadcast <- msg
	}

	g.Name = name
}

func (g *GameObject) GetName() string {
	g.RLock()
	defer g.RUnlock()

	return g.Name
}

func (g *GameObject) GetSymbol() rune {
	g.RLock()
	defer g.RUnlock()

	return g.Symbol
}

func (g *GameObject) SetColor(color tb.Attribute) {
	g.Lock()
	defer g.Unlock()

	g.Color = color
}

func (g *GameObject) GetColor() tb.Attribute {
	g.RLock()
	defer g.RUnlock()

	return g.Color
}

func (g *GameObject) JSON() ([]byte, error) {
	g.RLock()
	defer g.RUnlock()

	return json.Marshal(g)
}

func (g *GameObject) GetClass() string {
	g.RLock()
	defer g.RUnlock()

	return g.Class
}

func (g *GameObject) Event(action string) *z.Message {
	m := z.NewMessage(g.Class, g.ID, action)

	return m
}
