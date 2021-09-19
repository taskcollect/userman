package security

// TODO(codian) make this work with argon2id
// ! insecure !
func VerifySecrets(given string, stored string) bool {
	return given == stored
}
