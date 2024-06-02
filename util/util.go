package util

import (
	"math/rand"
	"strconv"
	"strings"
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

// ParseDurationString converts a duration string like "7d" to a time.Time value after that duration from now
func ParseDurationString(durationStr string) (time.Time, error) {
	now := time.Now()

	// If the string contains "d", handle it separately
	if strings.HasSuffix(durationStr, "d") {
		value, err := strconv.Atoi(durationStr[:len(durationStr)-1])
		if err != nil {
			return time.Time{}, err
		}
		return now.Add(time.Duration(value) * 24 * time.Hour), nil
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return time.Time{}, err
	}

	return now.Add(duration), nil
}
