package rsa

import (
	"math/big"
)

type PrivateKey struct {
	// Pointer to the public key which holds the modulus
	public *PublicKey

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
