package util

import (
	cryptoRand "crypto/rand"
	"fmt"
	"math/big"
	"time"
)

const VerificationOTPLifetime = 10 * time.Minute

// GenerateVerificationOTP returns a cryptographically secure six-digit OTP.
func GenerateVerificationOTP() (string, error) {
	otpNumber, err := cryptoRand.Int(
		cryptoRand.Reader,
		big.NewInt(1_000_000),
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", otpNumber.Int64()), nil
}
