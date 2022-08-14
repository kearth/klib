package web

import (
	"encoding/base64"
)

// Base64 Decode
func Base64Decode(s string) (string, error) {
	r, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(r), nil
}
