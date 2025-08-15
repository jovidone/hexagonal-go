package utils

import (
        "errors"
        "fmt"
        "github.com/golang-jwt/jwt/v5"
        "time"
)

var jwtKey = []byte("secret-key")

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
