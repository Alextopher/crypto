package myrsa

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

// Test saving and reading a public key
func TestSavePublicKey(t *testing.T) {
	// Keygen
	private := Keygen(1024)

	// Create a buffer to save the public key
	buf := bytes.Buffer{}

	// Save the public key
	private.Public.Save(&buf)

	// Create a new public key from the buffer
	public := ReadPublicKey(&buf)

	// Check the public key
	if public.n.Cmp(private.Public.n) != 0 {
		t.Error("Public key modulus mismatch")
	}

	if public.e.Cmp(private.Public.e) != 0 {
		t.Error("Public key exponent mismatch")
	}
}

// Test saving and reading a private key
func TestSavePrivateKey(t *testing.T) {
	// Keygen
	private := Keygen(1024)

	// Create a buffer to save the private key
	buf := bytes.Buffer{}

	// Save the private key
	private.Save(&buf)

	// Create a new private key from the buffer
	key := ReadPrivateKey(&buf)

	// Check the private key
	if key.Public.n.Cmp(private.Public.n) != 0 {
		t.Error("Private key modulus mismatch")
	}

	if key.Public.e.Cmp(private.Public.e) != 0 {
		t.Error("Private key exponent mismatch")
	}

	if key.d.Cmp(private.d) != 0 {
		t.Error("Private key private exponent mismatch")
	}

	if key.p.Cmp(private.p) != 0 {
		t.Error("Private key prime 1 mismatch")
	}

	if key.q.Cmp(private.q) != 0 {
		t.Error("Private key prime 2 mismatch")
	}
}

// Test keygen
func TestKeygen(t *testing.T) {
	keysizes := []int{128, 256, 512, 1024, 2048}

	for _, keysize := range keysizes {
		// Keygen
		private := Keygen(uint(keysize))

		// Check the prime sizes
		if private.p.BitLen() != keysize {
			t.Error("Prime 1 size mismatch", private.p.BitLen(), keysize)
		}

		if private.q.BitLen() != keysize {
			t.Error("Prime 2 size mismatch", private.q.BitLen(), keysize)
		}

		// Check that n is p * q
		if private.Public.n.Cmp(new(big.Int).Mul(private.p, private.q)) != 0 {
			t.Error("Modulus mismatch")
		}

		pm1 := new(big.Int).Sub(private.p, big.NewInt(1))
		qm1 := new(big.Int).Sub(private.q, big.NewInt(1))
		phi := new(big.Int).Mul(pm1, qm1)

		// Check that e is coprime to phi
		if gcd(private.Public.e, phi).Cmp(big.NewInt(1)) != 0 {
			t.Error("Public exponent not coprime to phi")
		}

		// d * e = 1 mod phi
		if new(big.Int).Mod(new(big.Int).Mul(private.d, private.Public.e), phi).Cmp(big.NewInt(1)) != 0 {
			t.Error("Private exponent not correct")
		}
	}
}

func TestEncrypDecrypt(t *testing.T) {
	// Keygen
	private := Keygen(128)

	// Create a buffer for the message
	msgbuf := bytes.Buffer{}

	// Write a message to the buffer
	msgbuf.WriteString("Hello World!")

	// Create a buffer for the encrypted message
	encbuf := bytes.Buffer{}

	// Encrypt the message
	private.Public.Encrypt(&msgbuf, &encbuf)

	// Create a buffer for the decrypted message
	decbuf := bytes.Buffer{}

	// Decrypt the message
	private.Decrypt(&encbuf, &decbuf)

	// Check the decrypted message
	if decbuf.String() != "Hello World!" {
		t.Error("Decrypted message does not match original")
	}
}
