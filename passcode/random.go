package passcode

import (
	"math"
	"math/rand"
	"strconv"
)

func times(str string, n int) (out string) {
	for i := 0; i < n; i++ {
		out += str
	}
	return
}

// Left left-pads the string with pad up to len runes
// len may be exceeded if
func padLeft(str string, length int, pad string) string {
	return times(pad, length-len(str)) + str
}

func Generate(length int) string {
	max := int(math.Pow(float64(10), float64(length))) - 1
	return padLeft(strconv.Itoa(rand.Intn(max)), length, "0")
}
