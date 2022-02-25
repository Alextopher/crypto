package main

import (
	"fmt"
	"os"
	"strconv"

	rsa "github.com/Alextopher/crypto/hw2/myrsa"
)

func main() {
	// Keygen usage ./rsa keygen <key size> <public_key> <private_key>
	// Encrypt usage ./rsa encrypt <public_key> <plaintext> <ciphertext>
	// Decrypt usage ./rsa decrypt <private_key> <ciphertext> <plaintext>
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./rsa <keygen|encrypt|decrypt>")
		fmt.Println("Usage: ./rsa keygen <key size> <public_key> <private_key>")
		fmt.Println("Usage: ./rsa encrypt <public_key> <plaintext> <ciphertext>")
		fmt.Println("Usage: ./rsa decrypt <private_key> <ciphertext> <plaintext>")
		os.Exit(1)
	}

	if os.Args[1] == "keygen" {
		if len(os.Args) != 5 {
			panic("Usage: ./rsa keygen <key size> <public_key> <private_key>")
		}

		// Parse the key size
		keySize, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Println("Invalid key size")
			os.Exit(1)
		}

		if keySize%8 != 0 {
			fmt.Println("Key size must be a multiple of 8")
			os.Exit(1)
		}

		private := rsa.Keygen(uint(keySize))

		// Open the public key file
		publicFile, err := os.Create(os.Args[3])
		if err != nil {
			fmt.Println("Error opening public key file", err)
			os.Exit(1)
		}
		defer publicFile.Close()

		// Open the private key file for writing (only user rewritable)
		privateFile, err := os.OpenFile(os.Args[4], os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println("Error opening private key file", err)
			os.Exit(1)
		}
		defer privateFile.Close()

		// Save the keys
		private.Public.Save(publicFile)
		private.Save(privateFile)
	} else if os.Args[1] == "encrypt" {
		if len(os.Args) != 5 {
			panic("Usage: ./rsa encrypt <public_key> <plaintext> <ciphertext>")
		}

		// Open the public key file
		publicFile, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println("Error opening public key file", err)
			os.Exit(1)
		}
		defer publicFile.Close()

		// Open the plaintext file
		plaintextFile, err := os.Open(os.Args[3])
		if err != nil {
			fmt.Println("Error opening plaintext file", err)
			os.Exit(1)
		}
		defer plaintextFile.Close()

		// Open the ciphertext file
		ciphertextFile, err := os.Create(os.Args[4])
		if err != nil {
			fmt.Println("Error opening ciphertext file", err)
			os.Exit(1)
		}
		defer ciphertextFile.Close()

		// Load the public key
		public := rsa.ReadPublicKey(publicFile)

		// Encrypt the plaintext
		public.Encrypt(plaintextFile, ciphertextFile)
	} else if os.Args[1] == "decrypt" {
		if len(os.Args) != 5 {
			panic("Usage: ./rsa decrypt <private_key> <ciphertext> <plaintext>")
		}

		// Open the private key file
		privateFile, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println("Error opening private key file", err)
			os.Exit(1)
		}
		defer privateFile.Close()

		// Open the ciphertext file
		ciphertextFile, err := os.Open(os.Args[3])
		if err != nil {
			fmt.Println("Error opening ciphertext file", err)
			os.Exit(1)
		}

		// Open the plaintext file
		plaintextFile, err := os.Create(os.Args[4])
		if err != nil {
			fmt.Println("Error opening plaintext file", err)
			os.Exit(1)
		}
		defer plaintextFile.Close()

		// Load the private key
		private := rsa.ReadPrivateKey(privateFile)

		// Decrypt the ciphertext
		private.Decrypt(ciphertextFile, plaintextFile)
	}
}
