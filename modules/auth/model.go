package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *LoginForm) HashPassword() {
	bytePassword := []byte(l.Password)
	sha256Bytes := sha256.Sum256(bytePassword)
	l.Password = hex.EncodeToString(sha256Bytes[:])
}
