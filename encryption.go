package main

import (
	"crypto/aes"
	"crypto/cipher"
	b64 "encoding/base64"
)

func createAesGCM(key string) cipher.AEAD {
	// create aes cipher
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	return aesGCM
}

func encrypt(email string) string {
	aesGCM := createAesGCM(config.CryptoKey)
	// encrypt email address
	encryptedEmail := make([]byte, aesGCM.NonceSize())
	ciphertext := aesGCM.Seal(encryptedEmail, encryptedEmail, []byte(email), nil)
	return b64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(encryptedEmail string) string {
	aesGCM := createAesGCM(config.CryptoKey)
	// decrypt email address
	sDec, err := b64.StdEncoding.DecodeString(encryptedEmail)
	if err != nil {
		print(err)
		return ""
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := sDec[:nonceSize], sDec[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		print(err)
		return ""
	}
	return string(plaintext)
}
