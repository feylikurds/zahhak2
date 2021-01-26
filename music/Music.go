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

package music

import (
	"time"

	zgo "../gameobjects"
)

type Music struct {
	zgo.GameObject

	sound *Sound

	background [][]int16
	sfx        map[string][][]int16

	bstream       chan []int16
	estream       chan [][]int16
	effects       chan [][]int16
	effectPlaying bool
}

func NewMusic() *Music {
	btoi16 := func(b []byte) (u []int16) {
		u = make([]int16, len(b)/2)

		for i, _ := range u {
			val := int16(b[i*2])
			val += int16(b[i*2+1]) << 8
			u[i] = val
		}

		return
	}

	chunkisize := func(data []int16, bufsize int) [][]int16 {
		chunk := int(bufsize)
		l := len(data)
		n := l / chunk
		chunks := make([][]int16, n)

		if n > 0 {
			for i := 0; i < n; i++ {
				chunks[i] = data[i*chunk : (i+1)*chunk]
			}

			size := len(chunks[n-1])

			if size < int(bufsize) {
				rest := make([]int16, int(bufsize)-size)
				chunks[n-1] = append(chunks[n-1], rest...)
			}
		} else {
			chunks = make([][]int16, 1)
			chunks[0] = data[:]
			size := len(chunks[0])

			if size < int(bufsize) {
				rest := make([]int16, int(bufsize)-size)
				chunks[0] = append(chunks[0], rest...)
			}
		}

		return chunks
	}

	backgroundInfo := ReadWavData("background.wav")

	sampleRate := int(backgroundInfo.SampleRate)
	numChannels := int(backgroundInfo.NumChannels)
	bytesPerFrame := numChannels * 2
	periodSize := int(float32(bytesPerFrame) * PERIOD_TIME)
	period := PERIOD
	bufferSize := period * periodSize
	latency := (bufferSize / (sampleRate * bytesPerFrame)) * 1000 * 1000

	backgroundData := chunkisize(btoi16(backgroundInfo.data), bufferSize)

	sfx := map[string][][]int16{}

	fireInfo := ReadWavData("fire.wav")
	sfx["fire"] = chunkisize(btoi16(fireInfo.data), bufferSize)

	teleportInfo := ReadWavData("teleport.wav")
	sfx["teleport"] = chunkisize(btoi16(teleportInfo.data), bufferSize)

	explodeInfo := ReadWavData("explode.wav")
	sfx["explode"] = chunkisize(btoi16(explodeInfo.data), bufferSize)

	itemInfo := ReadWavData("item.wav")
	sfx["item"] = chunkisize(btoi16(itemInfo.data), bufferSize)

	dieInfo := ReadWavData("die.wav")
	sfx["die"] = chunkisize(btoi16(dieInfo.data), bufferSize)

	sound := NewSound("default", sampleRate, numChannels, bytesPerFrame, periodSize, period, bufferSize, latency)

	return &Music{
		sound:      sound,
		background: backgroundData,
		sfx:        sfx,
		bstream:    make(chan []int16, 1),
		estream:    make(chan [][]int16, 1),
		effects:    make(chan [][]int16, 1),
	}
}

func (m *Music) Run() {
	if m.sound.err != nil {
		return
	}

	go func() {
		for {
			if !m.Running() {
				time.Sleep(200 * time.Millisecond)

				continue
			}

			<-m.sound.ready

			select {
			case effect := <-m.estream:
				m.effectPlaying = true
				m.sound.ready <- true

				m.sound.flush()

				for _, chunk := range effect {
					<-m.sound.ready
					m.sound.playback(chunk)
				}

				m.effectPlaying = false

			default:
				chunk := <-m.bstream
				m.sound.playback(chunk)
			}
		}
	}()

	go func() {
		for {
			effect := <-m.effects
			m.estream <- effect
		}
	}()
}

func (m *Music) Background() {
	go func() {
		for {
			for _, chunk := range m.background {
				m.bstream <- chunk
			}
		}
	}()
}

func (m *Music) Play(effect string) {
	if m.effectPlaying {
		return
	}

	select {
	case m.effects <- m.sfx[effect]:
	default:
		return
	}
}
