package domain

import "crypto/sha256"

func KeyFor(kind string, seed string) string {
	digest := sha256.Sum256([]byte(kind + ":" + seed))
	return kind + "-" + hexDigest(digest[:])
}

func hexDigest(input []byte) string {
	const alphabet = "0123456789abcdef"
	output := make([]byte, len(input)*2)
	for index, value := range input {
		output[index*2] = alphabet[value>>4]
		output[index*2+1] = alphabet[value&15]
	}
	return string(output)
}

func IsKeyFor(kind string, key string) bool {
	prefix := kind + "-"
	if len(key) <= len(prefix) {
		return false
	}
	return key[:len(prefix)] == prefix
}
