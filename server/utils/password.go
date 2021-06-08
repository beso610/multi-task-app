package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func NewSalt64() []byte {
	salt := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil
	}
	return salt
}

func HashPassword(password string, salt []byte) []byte {
	hashed := pbkdf2.Key([]byte(password), salt, 65536, 64, sha512.New)
	return hashed[:]
}

func ComparePassword(password string, salt string, hashedPassword string) (bool, error) {
	saltByte, err := base64.RawURLEncoding.DecodeString(salt)
	if err != nil {
		return false, err
	}

	hashedPasswordByte, err := base64.RawURLEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(hashedPasswordByte, HashPassword(password, saltByte)) == 1, nil
}
