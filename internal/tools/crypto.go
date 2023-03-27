// Package tools holds some methods useful for Gophkeeper App
package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// EncryptString is used to encrypt string data
func EncryptString(plaintext string, key string) (string, error) {
	plaintextBytes := []byte(plaintext)
	keyBytes := []byte(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintextBytes, nil)
	ciphertextString := base64.StdEncoding.EncodeToString(ciphertext)

	return ciphertextString, nil
}

// DecryptString is used to decrypt string data
func DecryptString(ciphertextString string, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextString)
	if err != nil {
		return "", err
	}

	keyBytes := []byte(key)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	plaintext := string(plaintextBytes)
	return plaintext, nil
}

// GenerateKey used to generate a size 16 key from a string
func GenerateKey(input string) string {
	hash := sha256.Sum256([]byte(input))
	key := hash[:16]

	return hex.EncodeToString(key)
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
