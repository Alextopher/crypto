package main

import (
	"io"
	"math/bits"
)

// A simple linear feedback shift register
// https://en.wikipedia.org/wiki/Linear-feedback_shift_register
type LFSR struct {
	state uint64
	taps  uint64
}

// Creates a new LFSR with the given state and position of the taps
func NewLFSR(start, taps uint64) *LFSR {
	return &LFSR{start, taps}
}

// Implements io.ByteReader
func (lsfr *LFSR) ReadByte() (byte, error) {
	var b byte = 0
	for i := 0; i < 8; i++ {
		// Pairty
		bit := uint64(bits.OnesCount64(lsfr.state&lsfr.taps)) & 1
		lsfr.state = bit<<63 | lsfr.state>>1
		b = b<<1 | byte(bit)
	}
	return b, nil
}

// Implements io.Reader
func (lsfr *LFSR) Read(p []byte) (n int, err error) {
	for i := range p {
		b, _ := lsfr.ReadByte()
		p[i] = b
	}
	return len(p), nil
}

// Mixes the given io.Reader with randomness from the LFSR
func (lsfr *LFSR) Mix(r io.Reader, w io.Writer) {
	// Allocate some buffers
	buf := make([]byte, 4096)
	rand := make([]byte, 4096)
	enc := make([]byte, 4096)

	for {
		// Read from the reader
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		buf = buf[:n]

		// Read equal amount of random bytes
		n, _ = lsfr.Read(rand[:n])

		// XOR the two
		for i := 0; i < n; i++ {
			enc[i] = buf[i] ^ rand[i]
		}

		// Write the encrypted bytes
		_, err = w.Write(enc[:n])

		if err != nil {
			panic(err)
		}
	}
}
