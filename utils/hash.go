package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(content string) string {
	hashedContent, _ := bcrypt.GenerateFromPassword([]byte(content), bcrypt.DefaultCost)
	return string(hashedContent)
}

func CompareHashed(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
