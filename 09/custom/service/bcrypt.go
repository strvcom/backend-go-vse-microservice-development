package service

import (
	"crypto/hmac"
	"encoding/base64"
	"os"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

var (
	pepper []byte
)

const (
	bcryptMaxPasswordLength = 72
)

func init() {
	p, ok := os.LookupEnv("pepper")
	if !ok {
		os.Exit(1)
	}
	pepper = []byte(p)
}

func HashPassword(password []byte) ([]byte, error) {
	passwordHash, err := encodedSHA512(password)
	if err != nil {
		return nil, err
	}
	return bcrypt.GenerateFromPassword(passwordHash[:bcryptMaxPasswordLength], bcrypt.DefaultCost)
}

func CompareHashAndPassword(hash, password []byte) bool {
	passwordHash, err := encodedSHA512(password)
	if err != nil {
		return false
	}
	if err = bcrypt.CompareHashAndPassword(hash, passwordHash[:bcryptMaxPasswordLength]); err != nil {
		return false
	}
	return true
}

func encodedSHA512(value []byte) ([]byte, error) {
	shaHasher := hmac.New(sha3.New512, pepper)
	if _, err := shaHasher.Write(value); err != nil {
		return nil, err
	}
	hash := shaHasher.Sum(nil)
	return []byte(base64.StdEncoding.EncodeToString(hash)), nil
}
