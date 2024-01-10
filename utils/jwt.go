package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	AccessTokenType  = "ACCESS"
	RefreshTokenType = "REFRESH"
)

var (
	ErrExpiredToken = errors.New("Token has expired")
	ErrInvalidToken = errors.New("Token is invalid")
)

type JWTClaims struct {
	Sub       uuid.UUID `json:"sub"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (c *JWTClaims) Valid() error {
	if time.Now().After(c.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

type Strategy interface {
	GenerateToken(userId uuid.UUID, email string) (string, string, error)
	VerifyToken(token, tokenType string) (*JWTClaims, error)
}

type JWTStrategy struct {
	accessTokenSecret  string
	refreshTokenSecret string
}

func NewJWTStrategy(accessTokenSecret, refreshTokenSecret string) Strategy {
	return &JWTStrategy{
		accessTokenSecret:  accessTokenSecret,
		refreshTokenSecret: refreshTokenSecret,
	}
}

func (strategy *JWTStrategy) GenerateToken(userId uuid.UUID, email string) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"sub":        userId,
		"email":      email,
		"issued_at":  time.Now(),
		"expired_at": time.Now().Add(15 * time.Minute),
	}

	refreshTokenClaims := jwt.MapClaims{
		"sub":        userId,
		"email":      email,
		"issued_at":  time.Now(),
		"expired_at": time.Now().Add(7 * 24 * time.Hour),
	}

	accessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessTokenClaims,
	).SignedString([]byte(strategy.accessTokenSecret))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshTokenClaims,
	).SignedString([]byte(strategy.refreshTokenSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (strategy *JWTStrategy) VerifyToken(token, tokenType string) (*JWTClaims, error) {
	tokenSecret := strategy.accessTokenSecret
	if tokenType == RefreshTokenType {
		tokenSecret = strategy.refreshTokenSecret
	}
	keyFunc := func(jwtToken *jwt.Token) (interface{}, error) {
		_, ok := jwtToken.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(tokenSecret), nil
	}

	accessToken, err := jwt.ParseWithClaims(token, &JWTClaims{}, keyFunc)

	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrInvalidToken) {
			return nil, ErrInvalidToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := accessToken.Claims.(*JWTClaims)

	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
