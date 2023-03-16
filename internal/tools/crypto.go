package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// EncryptString encrypts the plaintext using AES encryption with the secret key.
func EncryptString(plaintext string, secret string) (string, error) {
	key := []byte(secret)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintextBytes := []byte(plaintext)
	paddedPlaintext := make([]byte, len(plaintextBytes)+(aes.BlockSize-len(plaintextBytes)%aes.BlockSize))
	copy(paddedPlaintext, plaintextBytes)

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	encryptedData := append(iv, ciphertext...)
	encodedData := base64.StdEncoding.EncodeToString(encryptedData)
	return encodedData, nil
}

// DecryptString decrypts the ciphertext using AES decryption with the secret key.
func DecryptString(ciphertext string, secret string) (string, error) {
	key := []byte(secret)

	decodedData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	iv := decodedData[:aes.BlockSize]
	ciphertextBytes := decodedData[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	paddedPlaintext := make([]byte, len(ciphertextBytes))
	mode.CryptBlocks(paddedPlaintext, ciphertextBytes)

	plaintextBytes := make([]byte, len(paddedPlaintext))
	copy(plaintextBytes, paddedPlaintext[:len(paddedPlaintext)-int(paddedPlaintext[len(paddedPlaintext)-1])])

	return string(plaintextBytes), nil
}
