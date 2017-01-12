// coalesce/util/auth.go

package util

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"time"

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

func NiceTime(oldTime int64) string {
	curTime := time.Now().Unix()
	seconds := curTime - oldTime
	var elapsed string

	switch {
	// < 2 minutes
	case seconds < 120:
		elapsed = strconv.Itoa(int(seconds))
		return elapsed + " seconds ago"

	// < 2 hours
	case seconds < 7200:
		elapsed = strconv.Itoa(int(seconds / 60))
		return elapsed + " minutes ago"

	// < 2 days
	case seconds < 172800:
		elapsed = strconv.Itoa(int(seconds / 60 / 60))
		return elapsed + " hours ago"

	// < 2 months
	case seconds < 5256000:
		elapsed = strconv.Itoa(int(seconds / 60 / 60 / 24))
		return elapsed + " days ago"

	// < 2 years
	case seconds < 63072000:
		elapsed = strconv.Itoa(int(seconds / 60 / 60 / 24 / 30))
		return elapsed + " months ago"

	// 2 years +
	default:
		elapsed = strconv.Itoa(int(seconds / 60 / 60 / 24 / 30 / 12))
		return elapsed + " years ago"
	}
}
func CheckToken(salt, pw, token string) error {
	hash := []byte(token)
	password := []byte(salt + pw)
	return bcrypt.CompareHashAndPassword(hash, password)
}