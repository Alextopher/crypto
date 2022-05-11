package myrsa

import (
	"io"
	"math/big"
)

// ErrMessageTooLong
// Error returned when attempting to encrypt a message which is too large for the public key.
type ErrMessageTooLong struct{}

func (e ErrMessageTooLong) Error() string {
	return "message too long for RSA public key size"
}

func (public *PublicKey) encryptBlock(block []byte) (*big.Int, error) {
	m := new(big.Int).SetBytes(block)

	if m.Cmp(public.n) >= 0 {
		return nil, ErrMessageTooLong{}
	}

	return new(big.Int).Exp(m, public.e, public.n), nil
}

func (public *PublicKey) Encrypt(r io.Reader, w io.Writer) {
	// block size
	bs := public.n.BitLen() / 16

	for {
		// read a block
		block := make([]byte, bs)
		n, err := r.Read(block)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// encrypt it
		ciphertext, err := public.encryptBlock(block[:n])

		if err != nil {
			panic(err)
		}

		// write it to the output
		w.Write([]byte(ciphertext.String()))

		// Seperate with a newline
		w.Write([]byte("\n"))
	}
}

func (public *PublicKey) EncryptByByte(r io.Reader, w io.Writer) {
	// read buffer
	buf := make([]byte, 4096)

	for {
		// read a block
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// loop for each byte
		for i := 0; i < n; i++ {
			b := buf[i]

			// encrypt it
			ciphertext := new(big.Int).Exp(big.NewInt(int64(b)), public.e, public.n)

			// write it to the output
			w.Write([]byte(ciphertext.String()))

			// Seperate with a newline
			w.Write([]byte("\n"))
		}
	}
}
