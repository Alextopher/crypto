package myrsa

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"runtime"
	"time"
)

// For loop to find the gcd of two numbers
func gcd(a, b *big.Int) *big.Int {
	for b.Cmp(big.NewInt(0)) != 0 {
		a, b = b, new(big.Int).Mod(a, b)
	}

	return a
}

// Generates a uniform random number between a and b
func randomBetween(a, b *big.Int) (*big.Int, error) {
	// make sure a < b
	if a.Cmp(b) > 0 {
		a, b = b, a
	}

	// Calculate the range
	r := new(big.Int).Sub(b, a)

	// Generate a random number between 0 and r
	n, err := rand.Int(rand.Reader, r)

	// Add a to the number
	n.Add(n, a)

	return n, err
}

// Miller-Rabin primality test for n, with confidence k
// Implemented using psuedo code from
// https://en.wikipedia.org/wiki/Miller%E2%80%93Rabin_primality_test#Miller%E2%80%93Rabin_test
func millerRabin(n *big.Int, k uint) bool {
	// Pre calculate some important values
	one := big.NewInt(1)
	two := big.NewInt(2)
	nm1 := new(big.Int).Sub(n, one)
	nm2 := new(big.Int).Sub(n, two)

	// write n as 2^r·d + 1 with d odd
	d := new(big.Int).Set(nm1)
	r := uint(0)
	for d.Bit(0) == 0 {
		d.Rsh(d, 1)
		r++
	}

	// WitnessLoop
WitnessLoop:
	for i := uint(0); i < k; i++ {
		// pick a random integer a in the range [2, n − 2]
		a, _ := randomBetween(two, nm2)

		// x ← a^d mod n
		x := new(big.Int).Exp(a, d, n)

		// if x = 1 or x = n − 1 then
		if x.Cmp(one) == 0 || x.Cmp(nm1) == 0 {
			continue
		}

		// repeat r − 1 times:
		for j := uint(0); j < r-1; j++ {
			// x ← x^2 mod n
			x.Exp(x, two, n)

			// if x = n − 1 then
			if x.Cmp(nm1) == 0 {
				continue WitnessLoop
			}
		}

		// return “composite”
		return false
	}

	// return “probably prime”
	return true
}

// Tries to find a prime number of bits length. Returns nil if the attempt wasn't prime
func randomPrimeOneShot(bits uint, coprime *big.Int) *big.Int {
	// Generate random bytes
	b := make([]byte, bits/8)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// Set first and last bits to 1
	b[0] |= 0x80
	b[len(b)-1] |= 0x01

	// Create big int
	i := new(big.Int)
	i.SetBytes(b)

	// Make sure the prime is coprime to our e
	if coprime != nil {
		d := gcd(i, coprime)
		if d.Cmp(big.NewInt(1)) != 0 {
			return nil
		}
	}

	// Check if it is prime
	if millerRabin(i, WITNESS_COUNT) {
		return i
	}

	return nil
}

// Generates a random prime number of size bits that is coprime to e
// Sends the prime to the channel c when it is found
// Stops search when something is sent to the channel stop
func randomPrime(bits uint, e *big.Int, primes chan<- *big.Int, stop chan uint) {
	var count uint = 0

	for {
		select {
		case <-stop:
			stop <- count
			return
		default:
			prime := randomPrimeOneShot(bits, e)
			if prime != nil {
				select {
				case primes <- prime:
					<-stop
					stop <- count
				case <-stop:
					stop <- count
				}
			}
			count++
		}
	}
}

// RSA key genenerator
func Keygen(keySize uint) *PrivateKey {
	// Choose e
	e := big.NewInt(65537)

	// Generate p and q in parallel
	start := time.Now()
	primes := make(chan *big.Int)
	stops := make([]chan uint, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		stops[i] = make(chan uint)
		go randomPrime(keySize, e, primes, stops[i])
	}

	// Wait for the primes to be generated
	p := <-primes
	q := <-primes

	// Clean up the goroutines
	total := uint(0)
	for i := 0; i < runtime.NumCPU(); i++ {
		stops[i] <- 0
		total += <-stops[i]
	}

	fmt.Println("Generated p and q in", time.Since(start), "after", total, "tries")

	// Order p and q
	if p.Cmp(q) < 0 {
		p, q = q, p
	}

	// n = p * q
	n := new(big.Int).Mul(p, q)
	pm1 := new(big.Int).Sub(p, big.NewInt(1))
	qm1 := new(big.Int).Sub(q, big.NewInt(1))

	// Calculate phi(n) = (p-1) * (q-1)
	phi := new(big.Int).Mul(pm1, qm1)

	// Calculate d
	_, _, d := pulverizer(phi, e)
	d.Mod(d, phi)

	public := &PublicKey{
		n: n,
		e: e,
	}

	return &PrivateKey{public, d, p, q, nil, nil}
}
