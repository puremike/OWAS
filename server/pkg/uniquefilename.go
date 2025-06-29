package pkg

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateUniqueFileName() string {
	b := make([]byte, 12)

	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}
