package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var SecretKey = []byte("iamtheyammer")

type InvalidSignatureTypeError struct {
	TokenAlg jwt.SigningMethod
}

func (e InvalidSignatureTypeError) Error() string {
	return fmt.Sprintf("invalid jwt signing algorithm: %s specified but only hs256 is supported", e.TokenAlg)
}

func GenerateToken(userId int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["sub"] = userId

	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", fmt.Errorf("error generating signed jwt: %w", err)
	}

	return tokenString, nil
}

func ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, InvalidSignatureTypeError{TokenAlg: token.Method}
		}
		return SecretKey, nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing jwt: %w", err)
	}

	sub := token.Claims.(jwt.MapClaims)["sub"]

	return sub
}
