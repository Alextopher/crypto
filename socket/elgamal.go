package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Unlike AES ElGamal is a private/public key system
type ElGamalPrivateKey struct {
	// a
	a *big.Int

	// a pointer to the public key
	public *ElGamalPublicKey
}

type ElGamalPublicKey struct {
	// The prime p
	p *big.Int
	// The generator of the group
	g *big.Int
	// The public key, h = g^a mod p
	h *big.Int
}

// Geneates a new ElGamal key pair.
// Takes size of the prime as an argument in bits.
func Keygen(keysize int) (*ElGamalPrivateKey, *ElGamalPublicKey) {
	// generate a random prime
	p, err := rand.Prime(rand.Reader, keysize)
	if err != nil {
		panic("could not generate random prime")
	}

	// generate a random generator
	g, err := rand.Int(rand.Reader, p)
	if err != nil {
		panic("could not generate random generator")
	}

	// generate a random private key
	a, err := rand.Int(rand.Reader, p)
	if err != nil {
		panic("could not generate random private key")
	}

	// compute the public key
	h := new(big.Int).Exp(g, a, p)

	// create the private key
	private := &ElGamalPrivateKey{a, &ElGamalPublicKey{p, g, h}}

	return private, private.public
}

// This is the type that encodes messages to be sent over the wire
type ElGamalCipherText struct {
	// The shared secret
	shared *big.Int

	// The encrypted message
	ciphertext *big.Int

	// The number of bytes encoded in the ciphertext
	size int
}

// Encrypts a single message using ElGamal. The message must be less than the prime.
// Returns the shared secret followed by the encrypted message
func (pk *ElGamalPublicKey) _encrypt(m *big.Int) (ElGamalCipherText, error) {
	if m.Cmp(pk.p) >= 0 {
		return ElGamalCipherText{}, fmt.Errorf("message too large")
	}

	// choose a random big int less than p
	b, err := rand.Int(rand.Reader, pk.p)
	if err != nil {
		panic("could not generate random number")
	}

	// compute fullmask (g^a)^b mod p
	fullmask := new(big.Int).Exp(pk.h, b, pk.p)

	// compute shared secret g^b mod p
	shared := new(big.Int).Exp(pk.g, b, pk.p)

	// compute the ciphertext
	ciphertext := new(big.Int).Mul(fullmask, m)
	ciphertext.Mod(ciphertext, pk.p)

	return ElGamalCipherText{shared: shared, ciphertext: ciphertext}, nil
}

// Encrypts a full message using ElGamal. Depending on the size of the message multiple messages will be returned.
func (pk *ElGamalPublicKey) Encrypt(message []byte) ([]*ElGamalCipherText, error) {
	// calculate the block size
	bs := (pk.p.BitLen() / 8) - 1

	// check if the message is not divisible by the block size
	lastBlockSize := len(message) % bs

	// handle everything except the last block
	ciphers := make([]*ElGamalCipherText, 0)
	for i := 0; i < len(message)-lastBlockSize; i += bs {
		// create a new big int from the message
		m := new(big.Int).SetBytes(message[i : i+bs])

		// encrypt the message
		c, err := pk._encrypt(m)
		if err != nil {
			return nil, err
		}

		c.size = bs
		ciphers = append(ciphers, &c)
	}

	// handle the last block
	if lastBlockSize > 0 {
		// create a new big int from the message
		m := new(big.Int).SetBytes(message[len(message)-lastBlockSize:])

		// encrypt the message
		c, err := pk._encrypt(m)
		if err != nil {
			return nil, err
		}

		c.size = lastBlockSize
		ciphers = append(ciphers, &c)
	}

	return ciphers, nil
}

// Decrypts a single message using ElGamal.
func (sk *ElGamalPrivateKey) _decrypt(cipher *ElGamalCipherText) (*big.Int, error) {
	if cipher.shared.Cmp(sk.public.p) >= 0 {
		return nil, fmt.Errorf("shared secret too large")
	}

	if cipher.ciphertext.Cmp(sk.public.p) >= 0 {
		return nil, fmt.Errorf("ciphertext too large")
	}

	// full mask (g^b)^a mod p
	fullmask := new(big.Int).Exp(cipher.shared, sk.a, sk.public.p)

	// compute the modular inverse of the full mask
	fullmaskInv := fullmask.ModInverse(fullmask, sk.public.p)

	// compute the decrypted message
	plaintext := new(big.Int).Mul(cipher.ciphertext, fullmaskInv)
	plaintext.Mod(plaintext, sk.public.p)

	return plaintext, nil
}

// Decrypts a full message using ElGamal.
func (sk *ElGamalPrivateKey) Decrypt(ciphers []*ElGamalCipherText) ([]byte, error) {
	// decrypt each ciphertext
	plaintext := make([]byte, 0)
	for _, cipher := range ciphers {
		p, err := sk._decrypt(cipher)
		if err != nil {
			return nil, err
		}

		b := make([]byte, 0, cipher.size)

		// check if the plaintext is smaller than the block size
		missingBytes := cipher.size - len(p.Bytes())
		if missingBytes > 0 {
			for i := 0; i < missingBytes; i++ {
				b = append(b, 0)
			}
		} else if missingBytes < 0 {
			return nil, fmt.Errorf("ciphertext too large")
		}

		b = append(b, p.Bytes()...)
		plaintext = append(plaintext, b[:cipher.size]...)
	}

	return plaintext, nil
}
