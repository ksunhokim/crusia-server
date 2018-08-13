package server

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var (
	ErrInvalidBlockSize = errors.New("server: invalid block size")
)

func Encrypt(key, iv, payload []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := payload
	b = PKCS5Padding(b, aes.BlockSize, len(payload))

	cipherText := make([]byte, len(b))

	stream := cipher.NewCBCEncrypter(block, iv)
	stream.CryptBlocks(cipherText, b)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(key, iv []byte, payload string) ([]byte, error) {
	cipherText, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		return nil, ErrInvalidBlockSize
	}

	decrypted := make([]byte, len(cipherText))
	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(decrypted, cipherText)

	return PKCS5UnPadding(decrypted), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
