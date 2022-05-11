package main

import (
	"fmt"
	"math/rand"
	"os"
)

// This is the smallest known 64-bit LSFR which visits every state exactly once
const taps uint64 = 0x800000000000000D

func printUsage() {
	fmt.Println("Usage: ./lsfr keygen <keyfile>")
	fmt.Println("       ./lsfr mix <keyfile> <infile> <outfile>")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	switch os.Args[1] {
	case "keygen":
		if len(os.Args) != 3 {
			printUsage()
		}

		// Open key file
		f, err := os.Create(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Generate key
		key := rand.Uint64()

		// Write key to file
		fmt.Fprintf(f, "%d\n", key)
	case "mix":
		if len(os.Args) != 5 {
			fmt.Println("Usage: ./lsfr mix <keyfile> <infile> <outfile>")
			os.Exit(1)
		}

		// Open the key file
		f, err := os.Open(os.Args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		// Read the key
		var key uint64
		fmt.Fscanf(f, "%d\n", &key)

		// Open the input file
		r, err := os.Open(os.Args[3])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer r.Close()

		// Open the output file
		w, err := os.Create(os.Args[4])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer w.Close()

		// Create the LSFR
		lsfr := NewLSFR(key, taps)

		// Mix the files
		if lsfr.Mix(r, w) != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		printUsage()
	}
}
