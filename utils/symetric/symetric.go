package symetric

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"os"
	"strings"
)

const blockSize = 16

func Block(s string) cipher.Block {
	secret := make([]byte, blockSize)
	_, err := strings.NewReader(s).Read(secret)
	if err != nil {
		slog.Error("Error reading application secret", "error", err)
		os.Exit(1)
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		slog.Error("Error creating cipher block", "error", err)
		os.Exit(1)
	}

	return block
}

func Encrypt[T ~string | ~[]byte](block cipher.Block, text T) (string, error) {
	iv := make([]byte, blockSize)
	_, err := rand.Read(iv)
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cipherText := make([]byte, len(plainText))
	cipher.NewCFBEncrypter(block, iv).
		XORKeyStream(cipherText, plainText)

	encrypted := make([]byte, 0, len(iv)+len(cipherText))
	encrypted = append(encrypted, iv...)
	encrypted = append(encrypted, cipherText...)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt[T ~string | ~[]byte](block cipher.Block, text T) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return "", err
	}

	iv, cipherText := encrypted[:blockSize], encrypted[blockSize:]
	plainText := make([]byte, len(cipherText))
	cipher.NewCFBDecrypter(block, iv).
		XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
