package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
)

func GenerateHash(input string, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	data := []byte(input)
	h.Write(data)
	return h.Sum(nil)
}
