package main

import (
	"fmt"
	"math/big"
)

// y^2 = x^3 + ax + b (mod p)
type EllipicCurve struct {
	a *big.Int
	b *big.Int
	p *big.Int
}

type Point struct {
	x *big.Int
	y *big.Int
}

func (p *Point) String() string {
	return fmt.Sprintf("(%s, %s)", p.x, p.y)
}

func (p *Point) SetString(x, y string) (*Point, bool) {
	var ok bool

	p.x, ok = new(big.Int).SetString(x, 10)
	if !ok {
		return nil, false
	}
	p.y, ok = new(big.Int).SetString(y, 10)
	if !ok {
		return nil, false
	}

	return p, true
}

func (p1 *Point) Equals(p2 *Point) bool {
	return p1.x.Cmp(p2.x) == 0 && p1.y.Cmp(p2.y) == 0
}

func (e *EllipicCurve) String() string {
	return fmt.Sprintf("y^2 = x^3 + %sx + %s (mod %s)", e.a, e.b, e.p)
}

// IsOnCurve checks y^2 = x^3 + ax + b (mod p) for a given point
func (ec *EllipicCurve) IsOnCurve(p *Point) bool {
	// y^2 = x^3 + ax + b (mod p)
	lhs := new(big.Int).Mul(p.y, p.y)
	lhs = new(big.Int).Mod(lhs, ec.p)

	// x^3
	x3 := new(big.Int).Mul(p.x, p.x)
	x3 = new(big.Int).Mul(x3, p.x)

	// ax
	ax := new(big.Int).Mul(ec.a, p.x)

	// x^3 + ax + b
	rhs := new(big.Int).Add(x3, ax)
	rhs = new(big.Int).Add(rhs, ec.b)
	rhs = new(big.Int).Mod(rhs, ec.p)

	return lhs.Cmp(rhs) == 0
}

// Order returns the order of the point on the curve
func (ec *EllipicCurve) Order(g *Point) int {
	fmt.Println("1", g)
	p := ec.Double(g)

	order := 2
	for p.x.Cmp(g.x) != 0 || p.y.Cmp(g.y) != 0 {
		fmt.Println(order, p)
		p = ec.Add(p, g)
		order++
	}

	return order
}

// m = 3x^2 + a / 2y
// x3 = m^2 - 2x
// y3 = m(x - x3) - y
func (ec *EllipicCurve) Double(p *Point) *Point {
	Debug("Double(%s)\n", p)

	// x^2
	x2 := new(big.Int).Mul(p.x, p.x)

	// 2x^2 (mul by bitshift)
	top := new(big.Int).Lsh(x2, 1)

	// 3x^2 (add)
	top = new(big.Int).Add(top, x2)

	// 3x^2 + a
	top = new(big.Int).Add(top, ec.a)
	top = new(big.Int).Mod(top, ec.p)

	// 2y (mul by bitshift)
	bottom := new(big.Int).Lsh(p.y, 1)
	bottom = new(big.Int).Mod(bottom, ec.p)

	Debug("m = 3(%d)^2 + %d / 2*%d\nm = %d / %d\n", p.x, ec.a, p.y, top, bottom)

	// calculate the modular inverse of bottom
	_, _, inv := pulverizer(ec.p, bottom)
	Debug("m = %d * %d\n", top, inv)
	inv = new(big.Int).Mod(inv, ec.p)

	// m = 3x^2 + a / 2y
	m := new(big.Int).Mul(top, inv)
	m = new(big.Int).Mod(m, ec.p)
	Debug("m = %d * %d\nm = %d\n", top, inv, m)

	// x3 = m^2 - 2x
	x3 := new(big.Int).Mul(m, m)
	x3 = new(big.Int).Sub(x3, p.x)
	x3 = new(big.Int).Sub(x3, p.x)
	x3 = new(big.Int).Mod(x3, ec.p)

	Debug("x3 = m^2 - 2x = %d\n", x3)

	// y3 = m(x - x3) - y
	y3 := new(big.Int).Sub(p.x, x3)
	y3 = new(big.Int).Mul(m, y3)
	y3 = new(big.Int).Sub(y3, p.y)
	y3 = new(big.Int).Mod(y3, ec.p)

	Debug("y3 = m(x - x3) - y = %d\n", y3)
	Debug("result: %s\n", &Point{x3, y3})

	return &Point{x3, y3}
}

// m = (y2 - y1) / (x2 - x1)
// x3 = m^2 - x2 - x1
// y3 = m(x1 - x3) - y1
func (ec *EllipicCurve) Add(p1, p2 *Point) *Point {
	Debug("Add(%s, %s)\n", p1, p2)

	// p1 + p1 = 2p1
	if p1.x.Cmp(p2.x) == 0 && p1.y.Cmp(p2.y) == 0 {
		return ec.Double(p1)
	}

	// Point at infinity is the identity
	if p1.x.Cmp(big.NewInt(0)) == 0 && p1.y.Cmp(big.NewInt(0)) == 0 {
		return p2
	}

	if p2.x.Cmp(big.NewInt(0)) == 0 && p2.y.Cmp(big.NewInt(0)) == 0 {
		return p1
	}

	// mirrored points return the point at infinity
	if p1.x.Cmp(p2.x) == 0 {
		return &Point{big.NewInt(0), big.NewInt(0)}
	}

	// m = (y2 - y1) / (x2 - x1)
	dy := new(big.Int).Sub(p2.y, p1.y)
	dy = new(big.Int).Mod(dy, ec.p)

	dx := new(big.Int).Sub(p2.x, p1.x)
	dx = new(big.Int).Mod(dx, ec.p)

	Debug("m = (y2 - y1) / (x2 - x1)\nm = (%d - %d) / (%d - %d) \nm = %d / %d\n", p2.y, p1.y, p2.x, p1.x, dy, dx)

	_, _, dxInverse := pulverizer(ec.p, dx)
	dxInverse = new(big.Int).Mod(dxInverse, ec.p)
	Debug("m = %d * %d\n", dy, dxInverse)

	m := new(big.Int).Mul(dy, dxInverse)
	m = new(big.Int).Mod(m, ec.p)
	Debug("m = %d\n", m)

	// x3 = m^2 - x2 - x1
	x3 := new(big.Int).Mul(m, m)
	x3 = new(big.Int).Sub(x3, p2.x)
	x3 = new(big.Int).Sub(x3, p1.x)
	x3 = new(big.Int).Mod(x3, ec.p)
	Debug("x3 = m^2 - x2 - x1 = (%d)^2 - %d - %d = %d\n", m, p2.x, p1.x, x3)

	// y3 = m(x1 - x3) - y1
	y3 := new(big.Int).Sub(p1.x, x3)
	y3 = new(big.Int).Mul(m, y3)
	y3 = new(big.Int).Sub(y3, p1.y)
	y3 = new(big.Int).Mod(y3, ec.p)
	Debug("y3 = m(x1 - x3) - y1 = %d * (%d - %d) - %d = %d\n", m, p1.x, x3, p1.y, y3)

	Debug("result: %s\n", &Point{x3, y3})

	return &Point{x3, y3}
}

func (ec *EllipicCurve) Sub(p1, p2 *Point) *Point {
	return ec.Add(p1, &Point{p2.x, new(big.Int).Neg(p2.y)})
}

// ScalarMult returns k*P using the double and add method.
func (ec *EllipicCurve) ScalarMult(k *big.Int, p *Point) *Point {
	sum := &Point{big.NewInt(0), big.NewInt(0)}
	g := &Point{
		x: new(big.Int).Set(p.x),
		y: new(big.Int).Set(p.y),
	}

	for i := 0; i < k.BitLen(); i++ {
		if k.Bit(i) == 1 {
			sum = ec.Add(sum, g)
		}
		g = ec.Double(g)
	}

	return sum
}

// The pulverizer
func pulverizer(a, b *big.Int) (d, x, y *big.Int) {
	Debug("Modular inverse of %d mod %d\n", b, a)
	x1, y1, x2, y2 := big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1)
	for b.Cmp(big.NewInt(0)) != 0 {
		q, r := new(big.Int).DivMod(a, b, new(big.Int))
		Debug("a: %d b: %d q: %d r: %d x1: %d y1: %d x2: %d y2: %d\n", a, b, q, r, x1, y1, x2, y2)
		a, b = b, r
		x1, x2 = x2, new(big.Int).Sub(x1, new(big.Int).Mul(q, x2))
		y1, y2 = y2, new(big.Int).Sub(y1, new(big.Int).Mul(q, y2))
	}

	return a, x1, y1
}
