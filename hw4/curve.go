package main

import (
	"fmt"
	"math/big"
)

// // The pulverizer
// // gcd(a, b) = d = xa + yb mod a
// func pulverizer(a, b *big.Int) (d, x, y *big.Int) {
// 	x1, y1, x2, y2 := big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)
// 	for b.Cmp(big.NewInt(0)) != 0 {
// 		q, r := new(big.Int).DivMod(a, b, new(big.Int))
// 		// fmt.Println("a:", a, "b:", b, "q:", q, "r:", r, "x1:", x1, "y1:", y1, "x2:", x2, "y2:", y2)
// 		a, b = b, r
// 		x1, x2 = x2, new(big.Int).Sub(x1, new(big.Int).Mul(q, x2))
// 		y1, y2 = y2, new(big.Int).Sub(y1, new(big.Int).Mul(q, y2))
// 	}

// 	return a, x1, y1
// }

// y^2 = x^3 + ax + b (mod p)
type EllipticCurve struct {
	A, B, P, N *big.Int
	G          *Point
}

func NewEllipticCurve(a, b, p, n *big.Int, G *Point) *EllipticCurve {
	return &EllipticCurve{
		A: a,
		B: b,
		P: p,
		N: n,
		G: G,
	}
}

// Stringer for EllipticCurve
func (ec *EllipticCurve) String() string {
	return fmt.Sprintf("y^2 = x^3 + %dx + %d (mod %d)", ec.A, ec.B, ec.P)
}

// Verifies that the given point is on the curve
// By checking that y^2 = x^3 + ax + b (mod p)
func (ec *EllipticCurve) IsOnCurve(p *Point) bool {
	y2 := new(big.Int).Mul(p.Y, p.Y)
	y2 = y2.Mod(y2, ec.P)

	x3 := new(big.Int).Mul(p.X, p.X)
	x3 = x3.Mul(x3, p.X)
	x3 = x3.Mod(x3, ec.P)

	axb := new(big.Int).Mul(p.X, ec.A)
	axb = axb.Add(axb, ec.B)

	rsh := x3.Add(x3, axb)
	rsh = rsh.Mod(rsh, ec.P)

	return y2.Cmp(rsh) == 0
}

// Point addition
// m = (y2 - y1) / (x2 - x1)
// x3 = m^2 - x1 - x2
// y3 = m * (x1 - x3) - y1
func (ec *EllipticCurve) Add(p1, p2 *Point) *Point {
	if p1.X.Cmp(p2.X) == 0 && p1.Y.Cmp(p2.Y) == 0 {
		return ec.Double(p1)
	}

	if p1.X.Cmp(p2.X) == 0 {
		return &Point{big.NewInt(0), big.NewInt(0)}
	}

	dy := new(big.Int).Sub(p2.Y, p1.Y)
	dx := new(big.Int).Sub(p2.X, p1.X)
	// find the modular inverse of dx
	m := new(big.Int).ModInverse(dx, ec.P)
	m = m.Mul(m, dy)
	m = m.Mod(m, ec.P)

	x3 := new(big.Int).Mul(m, m)
	x3 = x3.Sub(x3, p1.X)
	x3 = x3.Sub(x3, p2.X)
	x3 = x3.Mod(x3, ec.P)

	x1MinusX3 := new(big.Int).Sub(p1.X, x3)
	y3 := new(big.Int).Mul(m, x1MinusX3)
	y3 = y3.Sub(y3, p1.Y)
	y3 = y3.Mod(y3, ec.P)

	return &Point{x3, y3}
}

// Point subtraction
func (ec *EllipticCurve) Sub(p1, p2 *Point) *Point {
	// Negate p2.y mod p
	p2y := new(big.Int).Neg(p2.Y)
	p2y = p2y.Mod(p2y, ec.P)

	return ec.Add(p1, &Point{p2.X, p2y})
}

// Point doubling
// m = 3x^2 + a / 2y
// x3 = m^2 - 2x
// y3 = m * (x - x3) - y
func (ec *EllipticCurve) Double(p *Point) *Point {
	threeX2a := new(big.Int).Mul(p.X, p.X)
	threeX2a = threeX2a.Mul(threeX2a, big.NewInt(3))
	threeX2a = threeX2a.Add(threeX2a, ec.A)
	threeX2a = threeX2a.Mod(threeX2a, ec.P)

	twoY := new(big.Int).Add(p.Y, p.Y)
	twoY = twoY.Mod(twoY, ec.P)
	// find the modular inverse of twoY
	m := new(big.Int).ModInverse(twoY, ec.P)
	m = m.Mul(m, threeX2a)
	m = m.Mod(m, ec.P)

	x3 := new(big.Int).Mul(m, m)
	x3 = x3.Sub(x3, p.X)
	x3 = x3.Sub(x3, p.X)
	x3 = x3.Mod(x3, ec.P)

	x1MinusX3 := new(big.Int).Sub(p.X, x3)
	y3 := new(big.Int).Mul(m, x1MinusX3)
	y3 = y3.Sub(y3, p.Y)
	y3 = y3.Mod(y3, ec.P)

	return &Point{x3, y3}
}

// Fast scalar multiplication
// only works if k > 2
func (ec *EllipticCurve) ScalarMult(p *Point, k *big.Int) *Point {
	if k.Cmp(big.NewInt(2)) == -1 {
		panic("ScalarMult k must be greater than 2")
	}

	// copy the point
	r := CopyPoint(p)
	// copy k
	k = new(big.Int).Set(k)
	k = k.Sub(k, big.NewInt(2))

	// "square and multiply" but since we're using additive notation it's really
	// "double and add"
	for k.Cmp(big.NewInt(0)) > 0 {
		r = ec.Double(r)
		if k.Bit(0) == 1 {
			r = ec.Add(r, p)
		}
		k = k.Rsh(k, 1)
	}

	return r
}

type Point struct {
	X, Y *big.Int
}

func CopyPoint(p *Point) *Point {
	return &Point{
		X: new(big.Int).Set(p.X),
		Y: new(big.Int).Set(p.Y),
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func PointEqual(p1, p2 *Point) bool {
	return p1.X.Cmp(p2.X) == 0 && p1.Y.Cmp(p2.Y) == 0
}
