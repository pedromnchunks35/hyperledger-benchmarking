package utils

import (
	"math/rand"
	"time"
)

// ? Generate a random string
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	//? This is all the possible chars that we can have
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@"
	//? This is to create a random number generated from a range of int values that will become from a timestamp which is int based
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	//? range the length of the string we want
	for i := range b {
		//? assign an char to that byte slice.. every char represents a number in the slice
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
