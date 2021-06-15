package utils

import "crypto/sha256"

type PasswordHelper struct {
	salt string
}

func (ph *PasswordHelper) HashPassword(password string, email string) [32]byte {
	hash := sha256.Sum256([]byte(ph.salt + password + email))
	return hash
}
