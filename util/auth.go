// coalesce/util/auth.go

package util

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt() (string, error) {
	r := make([]byte, 32)
	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}

	hash := sha512.Sum512(r)
	return hex.EncodeToString(hash[:]), nil
}

func ComputeToken(salt, pw string) (string, error) {
	password := []byte(salt + pw)
	token, err := bcrypt.GenerateFromPassword(password, 14)
	return string(token), err
}

func CheckToken(salt, pw, token string) error {
	hash := []byte(token)
	password := []byte(salt + pw)
	return bcrypt.CompareHashAndPassword(hash, password)
}
