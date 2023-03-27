// Package tools holds some methods useful for Gophkeeper App
package tools

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func EncryptString(plaintext, secret string) (string, error) {
	key := []byte(secret)
	plaintextBytes := []byte(plaintext)

	for i := 0; i < len(plaintextBytes); i++ {
		plaintextBytes[i] ^= key[i%len(key)]
	}

	encodedCiphertext := base64.StdEncoding.EncodeToString(plaintextBytes)
	return encodedCiphertext, nil
}

func DecryptString(ciphertext, secret string) (string, error) {
	key := []byte(secret)

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(decodedCiphertext); i++ {
		decodedCiphertext[i] ^= key[i%len(key)]
	}

	return string(decodedCiphertext), nil
}

// GeneratePasswordHash hashes a user password in terms of security
func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash returns true if password is hashed in a valid way
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
