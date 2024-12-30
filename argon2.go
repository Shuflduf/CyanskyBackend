package main

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func Argon2Encrypt(password, salt string) string {
	// Parameters for Argon2
	// salt := make([]byte, 16)
	// _, err := rand.Read(salt)
	// if err != nil {
	// 	return "", fmt.Errorf("failed to generate salt: %w", err)
	// }
	saltBytes := stringToByteArray16(salt)

	// Parameters for Argon2
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	// Generate the hash
	hash := argon2.IDKey([]byte(password), saltBytes, time, memory, threads, keyLen)

	// Encode the salt and hash to base64 for storage
	saltBase64 := base64.RawStdEncoding.EncodeToString(saltBytes)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	// Return the encoded hash
	return fmt.Sprintf("%s$%s", saltBase64, hashBase64)
}

func stringToByteArray16(s string) []byte {
	var byteArray []byte
	copy(byteArray[:], s)
	return byteArray
}
