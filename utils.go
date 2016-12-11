package govatar

import "math/rand"

// RandInt returns random integer
func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// RandString returns random element from slice of string
func RandString(slice []string) string {
	return slice[RandInt(0, len(slice))]
}
