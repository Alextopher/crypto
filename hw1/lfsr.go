package main

import (
	"io"
	"math/bits"
)

// LSFR - linear shift feedback register
// Implemented by reading the wikipedia article
// https://en.wikipedia.org/wiki/Linear-feedback_shift_register
type LSFR struct {
	state uint64
	taps  uint64
}

// NewLSFR Creates a new linear shift feedback register with the given state and position of the taps
// start - The initial state of the LSFR (the key)
// taps - Position of taps. This should generate a maximial length cycle (visits every state exactly once) but this isn't enforced
func NewLSFR(start, taps uint64) *LSFR {
	return &LSFR{start, taps}
}

// ReadByte steps the LSFR 8 times to generate a byte
func (lsfr *LSFR) ReadByte() (byte, error) {
	var b byte
	for i := 0; i < 8; i++ {
		bit := uint64(bits.OnesCount64(lsfr.state&lsfr.taps)) & 1 // parity of the bits
		lsfr.state = bit<<63 | lsfr.state>>1                      // shift the state
		b = b<<1 | byte(bit)                                      // append the bit to the byte
	}
	return b, nil
}

// Read reads len(p) bytes from the LSFR and returns the number of bytes read
func (lsfr *LSFR) Read(p []byte) (n int, err error) {
	for i := range p {
		b, _ := lsfr.ReadByte()
		p[i] = b
	}
	return len(p), nil
}

// Mix xors the given io.Reader with randomness generated from the LSFR
func (lsfr *LSFR) Mix(r io.Reader, w io.Writer) error {
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
			return err
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
			return err
		}
	}

	return nil
}
