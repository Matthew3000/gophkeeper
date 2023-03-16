package tools

import (
	"encoding/base64"
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
