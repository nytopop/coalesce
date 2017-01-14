// coalesce/util/auth.go

package util

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

func GenerateSalt() (string, error) {
	r := make([]byte, 32)
	_, err := rand.Read(r)
	if err != nil {
		return "", err
	}

	hash := sha512.Sum512(r)
	return base64.StdEncoding.EncodeToString(hash[:4]), nil
}

func ComputeToken(salt, pw string) (string, error) {
	password := []byte(salt + pw)
	token, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

func CheckToken(salt, pw, token string) error {
	hash, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return err
	}

	password := []byte(salt + pw)
	err = bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return err
	}
	return nil
}
