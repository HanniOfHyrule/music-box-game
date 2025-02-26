package crypto

import (
	"math/rand"
	"time"
)

var seededRand *rand.Rand = rand.New(rand.NewSource(int64(time.Now().UnixNano())))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomAlphaNumericString(length int) string {
	return StringWithCharset(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}
