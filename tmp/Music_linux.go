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

package main

/*
#cgo LDFLAGS: -lasound

#include <alsa/asoundlib.h>

void *cbuf[1];
*/
import "C"

import (
	"errors"
	"unsafe"
)

type Sound struct {
	pcm_handle *C.snd_pcm_t

	err   error
	ready chan bool
}

func NewSound(device string, sampleRate, numChannels, bytesPerFrame, periodSize, period, bufferSize, latency uint32) (sound *Sound) {
	newErr := func(function string, code C.int) error {
		return errors.New(function + ": " + C.GoString(C.snd_strerror(code)))
	}

	sound = &Sound{}

	code := C.snd_pcm_open(&sound.pcm_handle, C.CString(device), C.SND_PCM_STREAM_PLAYBACK, 0)

	if code < 0 {
		sound.err = newErr("snd_pcm_open", code)

		return
	}

	C.snd_pcm_drop(sound.pcm_handle)

	code = C.snd_pcm_set_params(
		sound.pcm_handle,
		C.SND_PCM_FORMAT_S16,
		C.SND_PCM_ACCESS_RW_INTERLEAVED,
		C.uint(numChannels),
		C.uint(sampleRate),
		1,
		C.uint(latency))

	if code < 0 {
		sound.err = newErr("snd_pcm_set_params", code)

		return
	}

	var params *C.snd_pcm_sw_params_t

	code = C.snd_pcm_sw_params_malloc(&params)

	if code < 0 {
		sound.err = newErr("snd_pcm_hw_params_malloc", code)

		return
	}

	code = C.snd_pcm_sw_params_current(sound.pcm_handle, params)

	if code < 0 {
		sound.err = newErr("snd_pcm_sw_params_current", code)

		return
	}

	code = C.snd_pcm_sw_params_set_stop_threshold(sound.pcm_handle, params, 0)

	if code < 0 {
		sound.err = newErr("snd_pcm_sw_params_set_stop_threshold", code)

		return
	}

	code = C.snd_pcm_sw_params_set_silence_threshold(sound.pcm_handle, params, 0)

	if code < 0 {
		sound.err = newErr("snd_pcm_sw_params_set_silence_threshold", code)

		return
	}

	var boundary C.snd_pcm_uframes_t

	code = C.snd_pcm_sw_params_get_boundary(params, &boundary)

	if code < 0 {
		sound.err = newErr("snd_pcm_sw_params_get_boundary", code)

		return
	}

	code = C.snd_pcm_sw_params_set_silence_size(sound.pcm_handle, params, boundary)

	if code < 0 {
		sound.err = newErr("snd_pcm_sw_params_set_silence_size", code)

		return
	}

	code = C.snd_pcm_sw_params(sound.pcm_handle, params)

	if code < 0 {
		sound.err = newErr("snd_pcm_sw_params", code)

		return
	}

	code = C.snd_pcm_prepare(sound.pcm_handle)

	if code < 0 {
		sound.err = newErr("snd_pcm_prepare", code)

		return
	}

	sound.ready = make(chan bool, 1)
	sound.ready <- true

	return
}

func (s *Sound) playback(buf []int16) {
	defer func() {
		s.ready <- true
	}()

	if s.err != nil {
		return
	}

	for i := 0; i < len(buf); {
		written := C.snd_pcm_writei(s.pcm_handle, unsafe.Pointer(&buf[i]), C.snd_pcm_uframes_t(len(buf)-i))
		if written <= 0 {
			code := C.snd_pcm_recover(s.pcm_handle, C.int(written), 1)

			if code != 0 {
				return
			}
		} else {
			i += int(written)
		}
	}
}
