package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

type LoginForm struct {
	Type     string `json:"type" form:"type"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	ApiKey   string `json:"apiKey" form:"apiKey"`
}

func (l *LoginForm) HashPassword() {
	bytePassword := []byte(l.Password)
	sha256Bytes := sha256.Sum256(bytePassword)
	l.Password = hex.EncodeToString(sha256Bytes[:])
}
