package src

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"

	"github.com/charmbracelet/log"
)

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func CreateHash(byte []byte) []byte {
	h := sha256.New()
	h.Write(byte)
	bs := h.Sum(nil)
	return bs
}

func VerifyHash(data []byte, expectedHash []byte) bool {
	hash := CreateHash(data)
	expect, err := hex.DecodeString(string(expectedHash))
	if err != nil {
		log.Error("Decode hash fail!")
	}
	return hmac.Equal(hash, expect)
}
