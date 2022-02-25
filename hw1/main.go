package main

import (
	"os"
)

const taps uint64 = 0x800000000000000D

func main() {
	// The key is the inital state of the tap
	var key uint64 = 12412341

	// I found a maximal period set of taps from this list (maybe tap choice could be part of the magic?)
	// https://users.ece.cmu.edu/~koopman/lfsr/64.txt
	lsfr1 := NewLFSR(key, taps)

	// encrypt message.txt to cipher.txt then decrypt cipher.txt to plain.txt
	msg, _ := os.Open("message.txt")
	cipher, _ := os.Create("cipher.txt")
	lsfr1.Mix(msg, cipher)

	cipher.Close()
	msg.Close()

	lsfr2 := NewLFSR(key, taps)

	// decrypt cipher.txt to plain.txt
	cipher, _ = os.Open("cipher.txt")
	plain, _ := os.Create("plain.txt")

	lsfr2.Mix(cipher, plain)
}
