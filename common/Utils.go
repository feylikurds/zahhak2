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
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"runtime"
	"sync"
)

var muRand = &sync.RWMutex{}

func UUID() string {
	muRand.Lock()
	defer muRand.Unlock()

	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)

	if n != len(uuid) || err != nil {
		panic(err)
	}

	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func RandomNumber(min int, max int) int {
	muRand.Lock()
	defer muRand.Unlock()

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))

	if err != nil {
		panic(err)
	}

	n := nBig.Int64()

	return int(n) + min
}

func MoveRandomly() (int, int) {
	x, y := 0, 0

	if flip() {
		if flip() {
			x += -1
		} else {
			x += 1
		}
	} else {
		if flip() {
			y += -1
		} else {
			y += 1
		}
	}

	return x, y
}

func MoveToGameObject(m IGameObject, p IGameObject) (int, int) {
	x, y := 0, 0
	mX, mY := m.GetPosition()
	pX, pY := p.GetPosition()
	deltaX, deltaY := mX-pX, mY-pY

	if deltaX < 0 {
		x += 1
	} else {
		x += -1
	}

	if deltaY < 0 {
		y += 1
	} else {
		y += -1
	}

	return x, y
}

func flip() bool {
	n := RandomNumber(1, 100)

	if n >= 50 {
		return true
	}

	return false
}

func Biased(bias int) bool {
	n := RandomNumber(1, 100)

	if n >= bias {
		return true
	}

	return false
}

func MouseToRelative(p ICreature, x, y int) (int, int) {
	pX, pY := p.GetPosition()
	deltaX, deltaY := pX-x, pY-y
	absDeltaX := math.Abs(float64(deltaX))
	absDeltaY := math.Abs(float64(deltaY))
	rX, rY := 0, 0

	if absDeltaX > absDeltaY {
		if deltaX < 0 {
			rX = 1
		} else {
			rX = -1
		}
	} else {
		if deltaY < 0 {
			rY = 1
		} else {
			rY = -1
		}
	}

	return rX, rY
}

func LogError(e error) {
	if e != nil {
		level := 1
		pc, fn, line, _ := runtime.Caller(level)
		file := runtime.FuncForPC(pc).Name()

		log.Printf("[error] in %s[%s:%d] %v", file, fn, line, e)
	}
}

func LogPanic(e error) {
	var fn, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])

	for _, pc := range pc[:n] {
		function := runtime.FuncForPC(pc)

		if function == nil {
			continue
		}

		file, line = function.FileLine(pc)
		fn = function.Name()

		log.Printf("[panic] in %s[%s:%d] %v", file, fn, line, e)
	}
}
