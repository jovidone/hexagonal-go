package utils

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"time"
)

var jwtKey []byte

// refreshTokenStore menyimpan refresh token yang valid untuk
// memungkinkan kontrol revokasi sederhana.
var refreshTokenStore = struct {
	sync.RWMutex
	tokens map[string]string
}{tokens: make(map[string]string)}

func init() {
	_ = godotenv.Load()
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString string) (string, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.UserID, nil
}

// RefreshClaims adalah klaim khusus untuk refresh token dengan masa berlaku
// yang lebih panjang.
type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateRefreshToken membuat refresh token baru dan menyimpannya di store.
func GenerateRefreshToken(userID string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &RefreshClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	refreshTokenStore.Lock()
	refreshTokenStore.tokens[tokenString] = userID
	refreshTokenStore.Unlock()

	return tokenString, nil
}

// ValidateRefreshToken memvalidasi refresh token dan memastikan token tersebut
// belum direvoke.
func ValidateRefreshToken(tokenString string) (string, error) {
	claims := &RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	refreshTokenStore.RLock()
	_, exists := refreshTokenStore.tokens[tokenString]
	refreshTokenStore.RUnlock()
	if !exists {
		return "", errors.New("token revoked")
	}

	return claims.UserID, nil
}

// RevokeRefreshToken menghapus refresh token dari store.
func RevokeRefreshToken(tokenString string) {
	refreshTokenStore.Lock()
	delete(refreshTokenStore.tokens, tokenString)
	refreshTokenStore.Unlock()
}
