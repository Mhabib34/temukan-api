package helper

import (
	"crypto/rand"
	"math/big"
)

func GenerateSecureOTP() (int64, error) {
	max := big.NewInt(900000) // range
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return n.Int64() + 100000, nil
}
