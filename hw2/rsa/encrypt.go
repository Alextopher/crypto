package rsa

import "math/big"

func (public *PublicKey) Encrypt(m *big.Int) *big.Int {
	return new(big.Int).Exp(m, public.e, public.n)
}
