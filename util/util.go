package util

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateUniqueCode() string {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	url := make([]byte, 10)
	for i := range url {
		url[i] = charset[r.Intn(len(charset))]
	}
	return string(url)
}
