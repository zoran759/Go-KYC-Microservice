package verification

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func generateIdempodencyKey() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, 21)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}