package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefhijklmnopqrstuv123456789"
const pwdalphabet = "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefhijklmnopqrstuv123456789.,!;"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int, alphabet string) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomPassword() string {
	return RandomString(13, pwdalphabet)
}

func RandomUsername() string {
	return RandomString(6, alphabet)
}

func RandomEmail() string {
	return RandomString(10, alphabet) + "@gmail.com"
}

func RandomDescription() string {
	return RandomString(30, pwdalphabet)
}
