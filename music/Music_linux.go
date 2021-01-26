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
// +build linux

package music

import (
	ga "github.com/cocoonlife/goalsa"
)

const PERIOD = 2
const PERIOD_TIME = 16

type Sound struct {
	pd *ga.PlaybackDevice

	err   error
	ready chan bool
}

func NewSound(device string, sampleRate, numChannels, bytesPerFrame, periodSize, period, bufferSize, latency int) (sound *Sound) {
	sound = &Sound{}

	sound.pd, sound.err = ga.NewPlaybackDevice("default", numChannels, ga.FormatS16LE, sampleRate,
		ga.BufferParams{
			BufferFrames: bufferSize})

	sound.ready = make(chan bool, 1)
	sound.ready <- true

	if sound.err != nil {
		panic(sound.err)
	}

	return
}

func (s *Sound) playback(buf []int16) {
	if s.err != nil {
		return
	}

	defer func() {
		s.ready <- true
	}()

	s.pd.Write(buf)
}

func (s *Sound) flush() {
}
