package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashString func is used to create hash based on
// string and const default is 14
func HashString(value string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	return string(bytes)
}

// CheckPasswordHash func is used to compare to hashes
func CheckPasswordHash(value, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(value))
	return err == nil
}
