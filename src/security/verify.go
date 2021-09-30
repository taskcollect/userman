package security

import "github.com/alexedwards/argon2id"

func Verify(given string, stored string) (bool, error) {
	return argon2id.ComparePasswordAndHash(given, stored)
}

func Hash(plaintext string) (string, error) {
	return argon2id.CreateHash(plaintext, argon2id.DefaultParams)
}
