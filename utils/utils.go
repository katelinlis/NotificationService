package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/avenir/notification-service/internal/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

func LoadPublicKey() ([]byte, error) {
	pubKeyPEM := []byte(os.Getenv("JWT_PUBLIC_KEY"))
	if len(pubKeyPEM) == 0 {
		const fallbackPath = "jwt_public.pem"
		data, err := os.ReadFile(fallbackPath)
		if err != nil {
			return nil, errors.New("JWT_PUBLIC_KEY is not set and failed to read from file: " + err.Error())
		}
		return data, nil
	}
	return pubKeyPEM, nil
}

func JWTParse(tokenString string) (MainClaims *model.MyCustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Load public key from environment variable
		pubKeyPEM, err := LoadPublicKey()
		if err != nil {
			return nil, err
		}
		// Parse PEM encoded public key
		block, _ := pem.Decode(pubKeyPEM)
		if block == nil {
			return nil, errors.New("failed to parse PEM block from JWT_PUBLIC_KEY")
		}
		parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaPubKey, ok := parsedKey.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("not RSA public key")
		}
		return rsaPubKey, nil
	})

	if err != nil {
		return MainClaims, err
	}

	if claims, ok := token.Claims.(*model.MyCustomClaims); ok {
		return claims, nil
	}
	return MainClaims, errors.New("invalid token")
}
