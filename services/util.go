package services

import "math/rand"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randCode(l int) string {
	b := make([]rune, l)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
