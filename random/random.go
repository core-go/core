package random

import "math/rand"

const arr = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Random(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = arr[rand.Intn(len(arr))]
	}
	return string(b)
}
