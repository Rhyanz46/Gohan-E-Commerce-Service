package utils

import (
	"math/rand"
	"time"
)

func RandomStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
