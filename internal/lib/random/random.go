package random

import (
	"math/rand"
	"time"
)

func RandomString(len int) string {
	rand := rand.New(rand.NewSource(time.Now().Unix()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	b := make([]rune, len)

	for i := range b {
		b[i] = chars[rand.Intn(len)]
	}

	return string(b)
}
