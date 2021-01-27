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
// +build windows

package music

/*
#cgo LDFLAGS: -lwinmm -I/usr/x86_64-w64-mingw32/include

#include <windows.h>
#include <mmsystem.h>

size_t wavehdrsize = sizeof(WAVEHDR);
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

const PERIOD = 1
const PERIOD_TIME = 0.25

type Sound struct {
	wfx      C.WAVEFORMATEX
	hwaveout C.HWAVEOUT
	wavehdr  *C.WAVEHDR

	err   error
	ready chan bool
}

func NewSound(device string, sampleRate, numChannels, bytesPerFrame, periodSize, period, bufferSize, latency int) (sound *Sound) {
	newErr := func(function string, code C.MMRESULT) error {
		var buf [1024]byte

		C.waveOutGetErrorText(code, C.LPSTR(unsafe.Pointer(&buf[0])), C.UINT(len(buf)))

		return errors.New(function + ": " + string(buf[:]))
	}

	sound = &Sound{}

	sound.wfx.wFormatTag = C.WAVE_FORMAT_PCM
	sound.wfx.nChannels = C.WORD(numChannels)
	sound.wfx.nSamplesPerSec = C.DWORD(sampleRate)
	sound.wfx.wBitsPerSample = 16
	sound.wfx.nBlockAlign = C.WORD(numChannels) * sound.wfx.wBitsPerSample / 8

	code := C.waveOutOpen(&sound.hwaveout, C.WAVE_MAPPER, &sound.wfx, 0, 0, C.CALLBACK_NULL)

	if code != C.MMSYSERR_NOERROR {
		sound.err = newErr("waveOutOpen", code)

		return
	}

	sound.ready = make(chan bool, 1)

	sound.ready <- true

	return
}

func (s *Sound) playback(buf []int16) {
	if s.err != nil {
		return
	}

	defer func() {
		s.ready <- true
	}()

	wavehdr := (*C.WAVEHDR)(C.malloc(C.wavehdrsize))
	C.memset(unsafe.Pointer(wavehdr), 0, C.wavehdrsize)
	wavehdr.lpData = C.LPSTR(unsafe.Pointer(&buf[0]))
	wavehdr.dwBufferLength = C.DWORD(len(buf) * 2)

	C.waveOutPrepareHeader(s.hwaveout, unsafe.Pointer(wavehdr), C.UINT(C.wavehdrsize))

	C.waveOutWrite(s.hwaveout, unsafe.Pointer(wavehdr), C.UINT(C.wavehdrsize))

	C.waveOutUnprepareHeader(s.hwaveout, unsafe.Pointer(wavehdr), C.UINT(C.wavehdrsize))

	s.wavehdr = wavehdr
}

func (s *Sound) flush() {
	if s.wavehdr != nil {
		var wavehdr C.WAVEHDR

		for s.wavehdr.dwFlags&C.WHDR_DONE == 0 {
			time.Sleep(time.Millisecond)
		}

		wdrsize := C.UINT(unsafe.Sizeof(wavehdr))
		C.waveOutUnprepareHeader(s.hwaveout, unsafe.Pointer(s.wavehdr), wdrsize)
		C.free(unsafe.Pointer(s.wavehdr))
	}
}

func (s *Sound) newErr(function string, code C.MMRESULT) error {
	var buf [1024]byte

	C.waveOutGetErrorText(code, C.LPSTR(unsafe.Pointer(&buf[0])), C.UINT(len(buf)))

	return errors.New(function + ": " + string(buf[:]))
}
