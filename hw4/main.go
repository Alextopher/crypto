package main

import (
	"fmt"
	"math/big"
)

const debug = true

func Debug(f string, v ...interface{}) {
	if debug {
		fmt.Printf(f, v...)
	}
}

func main() {
	a := big.NewInt(141)
	b := big.NewInt(162)
	p := big.NewInt(163)
	g := &Point{big.NewInt(6), big.NewInt(88)}
	n := big.NewInt(3)

	// n := big.NewInt(3)
	ec := &EllipicCurve{a, b, p}

	fmt.Println("Ellipic Curve:", ec)
	fmt.Println("G:", g, "isOnCurve:", ec.IsOnCurve(g))
	nG := ec.ScalarMult(n, g)
	fmt.Println("N*G:", nG, "isOnCurve:", ec.IsOnCurve(nG))

	// halfmask = (80, 7) = aG = 121G
	alice := big.NewInt(121)
	h := ec.ScalarMult(alice, g)

	fmt.Println("Halfmask:", h, "isOnCurve:", ec.IsOnCurve(h))

	// full mask = n * h = naG
	// f := ec.ScalarMult(n, h)
	f := ec.Add(nG, h)
	fmt.Println("Fullmask:", f, "isOnCurve:", ec.IsOnCurve(f))

	// cipher = (88, 71) = naG + m
	c := &Point{big.NewInt(88), big.NewInt(71)}

	// m = C - naG
	plain := ec.Sub(c, f)

	// plain text = cipher - fullmask
	fmt.Println("Plaintext:", plain, "isOnCurve:", ec.IsOnCurve(plain))

	// Find the period of G
	fmt.Println("Period of G:", ec.Order(g))
}
