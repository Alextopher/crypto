package myrsa

import (
	"bufio"
	"io"
	"math/big"
)

// Prepares Chinese remainder theorem
func (private *PrivateKey) precompute() {
	// c = m1 * p + m2 * q
	_, m1, m2 := pulverizer(private.p, private.q)
	m1 = m1.Mod(m1, private.p)
	m2 = m2.Mod(m2, private.p)

	private.ap = m1.Mul(m1, private.p)
	private.aq = m2.Mul(m2, private.q)
}

// Requires p > q
func (private *PrivateKey) decrypt_block(c *big.Int) *big.Int {
	if private.ap == nil || private.aq == nil {
		private.precompute()
	}

	// Fermat's little theorem
	c1, c2 := flt(c, private.d, private.p, private.q)

	// Applies the Chinese remainder theorem
	// m = c1 * p + c2 * q mod n
	m := new(big.Int).Mul(c1, private.ap)
	m.Add(m, new(big.Int).Mul(c2, private.aq))
	return m.Mod(m, private.Public.n)
}

func (private *PrivateKey) Decrypt(r io.Reader, w io.Writer) {
	// Each block is a hexadecimal number separated by a newline
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		chipher, _ := new(big.Int).SetString(scanner.Text(), 10)

		// Decrypt the block
		plaintext := private.decrypt_block(chipher)

		// Write the plaintext to the output
		w.Write(plaintext.Bytes())
	}
}
