package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestElGamalHelloWorld(t *testing.T) {
	message := []byte("1234")

	// Generate a new ElGamal key pair
	private, public := Keygen(32)

	// Encrypt the message
	ciphertext, err := public.Encrypt(message)
	if err != nil {
		t.Errorf("Could not encrypt message: %s", err)
	}

	// Decrypt the message
	plaintext, err := private.Decrypt(ciphertext)
	if err != nil {
		t.Errorf("Could not decrypt message: %s", err)
	}

	if string(plaintext) != string(message) {
		t.Errorf("Decrypted message does not match original message")
	}
}

// Parrallel Test to test large numbers of messages and keysizes
func TestElGamalParallel(t *testing.T) {
	var sizes = []int{1, 10, 100, 1000, 10000}
	message := make([]byte, 10000)
	rand.Read(message)

	// Test 128, 256, 512, 1024, 2048 keys
	for i := 128; i <= 2048; i *= 2 {
		// subtest for each key size
		t.Run(fmt.Sprintf("Key Size: %d", i), func(t *testing.T) {
			// Generate a new ElGamal key pair
			private, public := Keygen(i)

			// Test messages of different sizes
			for _, size := range sizes {
				ciphertext, err := public.Encrypt(message[:size])
				if err != nil {
					t.Errorf("Could not encrypt message: %s", err)
				}

				// Decrypt the message
				plaintext, err := private.Decrypt(ciphertext)
				if err != nil {
					t.Errorf("Could not decrypt message: %s", err)
				}

				if string(plaintext) != string(message[:size]) {
					t.Errorf("Decrypted message does not match original message")
				}
			}
		})
	}
}
