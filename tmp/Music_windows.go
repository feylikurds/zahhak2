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
// +build windows

package music

/*
#cgo LDFLAGS: -lwinmm -I/usr/x86_64-w64-mingw32/include

#include <windows.h>

size_t wavehdrsize = sizeof(WAVEHDR);
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

type Sound struct {
	hwaveout C.HWAVEOUT
	wavehdr  *C.WAVEHDR

	err   error
	ready chan bool
}

func NewSound(device string, sampleRate, numChannels, bytesPerFrame, periodSize, period, bufferSize, latency uint32) (sound *Sound) {
	newErr := func(function string, code C.MMRESULT) error {
		var buf [1024]byte

		C.waveOutGetErrorText(code, C.LPSTR(unsafe.Pointer(&buf[0])), C.UINT(len(buf)))

		return errors.New(function + ": " + string(buf[:]))
	}

	sound = &Sound{}

	var wavefx C.WAVEFORMATEX

	wavefx.wFormatTag = C.WAVE_FORMAT_PCM
	wavefx.wBitsPerSample = 16
	wavefx.nChannels = C.WORD(numChannels)
	wavefx.nBlockAlign = (wavefx.wBitsPerSample * wavefx.nChannels) / 8
	wavefx.nSamplesPerSec = C.DWORD(sampleRate)
	wavefx.nAvgBytesPerSec = wavefx.nSamplesPerSec * C.DWORD(wavefx.nBlockAlign)

	code := C.waveOutOpen(&sound.hwaveout, C.WAVE_MAPPER, &wavefx, 0, 0, C.CALLBACK_NULL)

	if code != C.MMSYSERR_NOERROR {
		sound.err = newErr("waveOutOpen", code)

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

	wavehdr := (*C.WAVEHDR)(C.malloc(C.wavehdrsize))
	C.memset(unsafe.Pointer(wavehdr), 0, C.wavehdrsize)
	wavehdr.lpData = C.LPSTR(unsafe.Pointer(&buf[0]))
	wavehdr.dwBufferLength = C.DWORD(len(buf))

	code := C.waveOutPrepareHeader(s.hwaveout, unsafe.Pointer(wavehdr), C.UINT(C.wavehdrsize))

	if code != C.MMSYSERR_NOERROR {
		s.err = s.newErr("waveOutPrepareHeader", code)

		return
	}

	if s.wavehdr != nil {
		for s.wavehdr.dwFlags&C.WHDR_DONE == 0 {
			time.Sleep(time.Millisecond)
		}
	}

	code = C.waveOutWrite(s.hwaveout, unsafe.Pointer(wavehdr), C.UINT(C.wavehdrsize))

	if code != C.MMSYSERR_NOERROR {
		s.err = s.newErr("waveOutWrite", code)

		return
	}

	for wavehdr.dwFlags&C.WHDR_DONE == 0 {
		time.Sleep(time.Millisecond)
	}

	code = C.waveOutUnprepareHeader(s.hwaveout, unsafe.Pointer(wavehdr), C.UINT(C.wavehdrsize))

	if code != C.MMSYSERR_NOERROR {
		s.err = s.newErr("waveOutUnprepareHeader", code)

		return
	}

	s.wavehdr = wavehdr
}

func (s *Sound) newErr(function string, code C.MMRESULT) error {
	var buf [1024]byte

	C.waveOutGetErrorText(code, C.LPSTR(unsafe.Pointer(&buf[0])), C.UINT(len(buf)))

	return errors.New(function + ": " + string(buf[:]))
}
