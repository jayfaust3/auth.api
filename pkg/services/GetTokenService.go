package services

import (
	"github.com/golang-jwt/jwt/v4"
)

func GetToken(encodedToken string) (string, error) {
	decodedToken, pemKey, err := ValidateToken(encodedToken)

	if err != nil {
		return encodedToken, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, decodedToken)

	byteKey := []byte(pemKey)
	reEncodedToken, err := token.SignedString(byteKey)

	if err != nil {
		return encodedToken, err
	}

	return reEncodedToken, err
}
