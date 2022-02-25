package myrsa

import (
	"bufio"
	"io"
	"math/big"
)

type PrivateKey struct {
	// Pointer to the public key which holds the modulus
	Public *PublicKey

	// The private exponent
	d *big.Int

	// The primes
	p, q *big.Int

	// Precomputed values for faster decryption
	ap *big.Int
	aq *big.Int
}

type PublicKey struct {
	// The modulus
	n *big.Int

	// The public exponent
	e *big.Int
}

// Save the private key to the writer
func (private *PrivateKey) Save(w io.Writer) {
	// Save the public key first
	private.Public.Save(w)

	// Save the private exponent
	w.Write([]byte(private.d.String()))
	w.Write([]byte("\n"))

	// Save the primes
	w.Write([]byte(private.p.String()))
	w.Write([]byte("\n"))

	w.Write([]byte(private.q.String()))
	w.Write([]byte("\n"))
}

func (public *PublicKey) Save(w io.Writer) {
	// First line is the modulus
	w.Write([]byte(public.n.String()))
	w.Write([]byte("\n"))

	// Second line is the public exponent
	w.Write([]byte(public.e.String()))
	w.Write([]byte("\n"))
}

func ReadPublicKey(r io.Reader) *PublicKey {
	scanner := bufio.NewScanner(r)

	// First line is the modulus
	scanner.Scan()
	n, _ := new(big.Int).SetString(scanner.Text(), 10)

	// Second line is the public exponent
	scanner.Scan()
	e, _ := new(big.Int).SetString(scanner.Text(), 10)

	return &PublicKey{n, e}
}

func ReadPrivateKey(r io.Reader) *PrivateKey {
	scanner := bufio.NewScanner(r)

	// First line is the modulus
	scanner.Scan()
	n, _ := new(big.Int).SetString(scanner.Text(), 10)

	// Second line is the public exponent
	scanner.Scan()
	e, _ := new(big.Int).SetString(scanner.Text(), 10)

	// Second line is the private exponent
	scanner.Scan()
	d, _ := new(big.Int).SetString(scanner.Text(), 10)

	// Third line is the p
	scanner.Scan()
	p, _ := new(big.Int).SetString(scanner.Text(), 10)

	// Fourth line is the q
	scanner.Scan()
	q, _ := new(big.Int).SetString(scanner.Text(), 10)

	return &PrivateKey{&PublicKey{n, e}, d, p, q, nil, nil}
}
