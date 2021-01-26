/*
Copyright (c) 2013 Tony Worm. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Tony Worm nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	"bufio"
	bin "encoding/binary"
	"fmt"
	"os"
)

type WavData struct {
	bChunkID  [4]byte // B
	ChunkSize uint32  // L
	bFormat   [4]byte // B

	bSubchunk1ID  [4]byte // B
	Subchunk1Size uint32  // L
	AudioFormat   uint16  // L
	NumChannels   uint16  // L
	SampleRate    uint32  // L
	ByteRate      uint32  // L
	BlockAlign    uint16  // L
	BitsPerSample uint16  // L

	bSubchunk2ID  [4]byte // B
	Subchunk2Size uint32  // L
	data          []byte  // L
}

func ReadWavData(fn string) (wav WavData) {
	ftotal, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("Error opening\n")
	}
	file := bufio.NewReader(ftotal)

	bin.Read(file, bin.BigEndian, &wav.bChunkID)
	bin.Read(file, bin.LittleEndian, &wav.ChunkSize)
	bin.Read(file, bin.BigEndian, &wav.bFormat)

	bin.Read(file, bin.BigEndian, &wav.bSubchunk1ID)
	bin.Read(file, bin.LittleEndian, &wav.Subchunk1Size)
	bin.Read(file, bin.LittleEndian, &wav.AudioFormat)
	bin.Read(file, bin.LittleEndian, &wav.NumChannels)
	bin.Read(file, bin.LittleEndian, &wav.SampleRate)
	bin.Read(file, bin.LittleEndian, &wav.ByteRate)
	bin.Read(file, bin.LittleEndian, &wav.BlockAlign)
	bin.Read(file, bin.LittleEndian, &wav.BitsPerSample)

	bin.Read(file, bin.BigEndian, &wav.bSubchunk2ID)
	bin.Read(file, bin.LittleEndian, &wav.Subchunk2Size)

	wav.data = make([]byte, wav.Subchunk2Size)
	bin.Read(file, bin.LittleEndian, &wav.data)

	return
}
