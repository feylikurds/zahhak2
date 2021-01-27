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
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	tb "github.com/nsf/termbox-go"
	ac "github.com/shiena/ansicolor"

	z "./common"
)

var game *Game
var quit bool
var prevX, prevY int

var muQuit = &sync.RWMutex{}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	f, e := os.OpenFile("zahhak2.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if e != nil {
		panic(e)
	}

	log.SetOutput(f)

	err := tb.Init()

	if err != nil {
		panic(err)
	}

	tb.SetOutputMode(tb.OutputNormal)
	tb.SetInputMode(tb.InputEsc | tb.InputMouse)
	tb.HideCursor()
}

func main() {
	defer mainRecover()

	quit = false

	//menu := zm.NewMenu()
	//menu.Options()

	config := z.NewConfig()
	//config.Multiplayer = menu.Multiplayer
	//config.Server = menu.Server
	//config.Host = menu.Host
	//config.Port = menu.Port
	//config.Name = menu.Name

	//if !menu.Quit {
	game = NewGame(config)
	game.Start()
	game.Play()

	go input()
	play()
	//}

	tb.Clear(tb.ColorDefault, tb.ColorDefault)
	tb.Close()

	w := ac.NewAnsiColorWriter(os.Stdout)
	text := "%s%s%s" + GOODBYE + "%s%s%s\n"

	fmt.Fprintf(w, text, "\x1b[31m", "\x1b[1m", "\x1b[40m", "\x1b[39m", "\x1b[49m", "\x1b[0m")
}

func play() {
	for !quit {
		game.Display()

		time.Sleep(100 * time.Millisecond)
	}
}

func input() {
	defer mainRecover()

	for !getQuit() {
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			switch ev.Key {
			case tb.KeyArrowUp:
				moveKey(0, -1)

			case tb.KeyArrowDown:
				moveKey(0, 1)

			case tb.KeyArrowRight:
				moveKey(1, 0)

			case tb.KeyArrowLeft:
				moveKey(-1, 0)

			case tb.KeySpace:
				shoot()

			case tb.KeyEnter:
				pause()

			case tb.KeyEsc:
				setQuit(true)

			default:
				ch := ev.Ch

				if ch == 'q' {
					setQuit(true)
				}
			}
		case tb.EventMouse:
			if ev.Key == tb.MouseLeft {
				moveMouse(ev.MouseX, ev.MouseY)
			} else if ev.Key == tb.MouseRight {
				shoot()
			}
		case tb.EventInterrupt:
			setQuit(true)

			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func getQuit() bool {
	muQuit.RLock()
	defer muQuit.RUnlock()

	return quit
}

func setQuit(newQuit bool) {
	muQuit.Lock()
	defer muQuit.Unlock()

	quit = newQuit
}

func moveKey(x, y int) {
	if x == prevX && y == prevY {
		return
	}

	game.MoveKey(x, y)

	prevX, prevY = x, y
}

func moveMouse(x, y int) {
	game.MoveMouse(x, y)
}

func shoot() {
	game.Fire()
}

func pause() {
	game.Pause()
}

func mainRecover() {
	if r := recover(); r != nil {
		e, ok := r.(error)

		if !ok {
			e = fmt.Errorf("%v", r)
			e = errors.New("Main: " + e.Error())
			z.LogPanic(e)
		} else {
			e = errors.New("Main: " + e.Error())
			z.LogPanic(e)
		}
	}
}
