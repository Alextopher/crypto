package myrsa

import (
	"fmt"
	"math/big"
)

const WITNESS_COUNT = 20

// The pulverizer
func pulverizer(a, b *big.Int) (d, x, y *big.Int) {
	x1, y1, x2, y2 := big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)
	for b.Cmp(big.NewInt(0)) != 0 {
		q, r := new(big.Int).DivMod(a, b, new(big.Int))
		fmt.Println("a:", a, "b:", b, "q:", q, "r:", r, "x1:", x1, "y1:", y1, "x2:", x2, "y2:", y2)
		a, b = b, r
		x1, x2 = x2, new(big.Int).Sub(x1, new(big.Int).Mul(q, x2))
		y1, y2 = y2, new(big.Int).Sub(y1, new(big.Int).Mul(q, y2))
	}

	return a, x1, y1
}

// Prepares a^b mod n as partial results to be used by the chinese remainder theorem
// x = (a mod p) ^ (b mod p-1) mod p
// y = (a mod q) ^ (b mod q-1) mod q
func flt(a, b, p, q *big.Int) (*big.Int, *big.Int) {
	x := make(chan *big.Int)
	y := make(chan *big.Int)

	go func() {
		// (a mod p) ^ (b mod p-1) mod p
		pm1 := new(big.Int).Sub(p, big.NewInt(1))
		x <- new(big.Int).Exp(new(big.Int).Mod(a, p), new(big.Int).Mod(b, pm1), p)
	}()

	go func() {
		// (a mod q) ^ (b mod q-1) mod q
		qm1 := new(big.Int).Sub(q, big.NewInt(1))
		y <- new(big.Int).Exp(new(big.Int).Mod(a, q), new(big.Int).Mod(b, qm1), q)
	}()

	return <-x, <-y
}
