package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessClaims struct {
	UID    uint   `json:"uid"`
	Role   string `json:"role"`
	UserId uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UID    uint `json:"uid"`
	UserId uint `json:"user_id"`
	// jti disimpan di RegisteredClaims.ID
	jwt.RegisteredClaims
}

func jwtSecret() ([]byte, error) {
	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		return nil, errors.New("AUTH_JWT_SECRET belum di-set")
	}
	return []byte(secret), nil
}

func GenerateAccessToken(uid uint, role string, ttl time.Duration) (token string, expiresIn int64, err error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", 0, err
	}
	now := time.Now()
	exp := now.Add(ttl)

	claims := AccessClaims{
		UID:    uid,
		Role:   role,
		UserId: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			// Opsional: isi Issuer/Audience jika perlu
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(secret)
	if err != nil {
		return "", 0, err
	}
	return signed, int64(time.Until(exp).Seconds()), nil
}

func GenerateRefreshToken(uid uint, ttl time.Duration) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	now := time.Now()
	exp := now.Add(ttl)

	jti, err := randomJTI(16)
	if err != nil {
		return "", err
	}

	claims := RefreshClaims{
		UID:    uid,
		UserId: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "refresh",
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

func randomJTI(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
