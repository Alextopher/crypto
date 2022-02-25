package main

import (
	"os"

	"github.com/Alextopher/crypto/hw2/rsa"
)

func main() {
	private := rsa.Keygen(2048)
	public := private.Public

	// Open message file
	msg, err := os.Open("message.txt")
	if err != nil {
		panic(err)
	}

	// Open ciphertext file
	cipher, err := os.Create("cipher.txt")
	if err != nil {
		panic(err)
	}

	// Encrypt message
	public.Encrypt(msg, cipher)

	msg.Close()
	cipher.Close()

	// Open ciphertext file
	cipher, err = os.Open("cipher.txt")
	if err != nil {
		panic(err)
	}

	// Open decrypted message file
	decrypted, err := os.Create("decrypted.txt")
	if err != nil {
		panic(err)
	}

	// Decrypt ciphertext
	private.Decrypt(cipher, decrypted)
}
