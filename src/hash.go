package src

import (
	"crypto/hmac"
	"crypto/sha256"
	// "encoding/hex"
)

func CreateHash(byte []byte) []byte {
	h := sha256.New()
	h.Write(byte)
	bs := h.Sum(nil)
	return bs
}

// func CreateHashString(byte []byte) string {
// 	return hex.EncodeToString(CreateHash(byte))
// }

func VerifyHash(data []byte, expectedHash []byte) bool {
	hash := CreateHash(data)
	return hmac.Equal(hash, expectedHash)
}

// func VerifyHashString(data []byte, expectedHash string) bool {
// 	hash := CreateHash(data)
// 	return hmac.Equal(hash, expectedHash)
// }
