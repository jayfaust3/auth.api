package services

import (
	"log"

	"github.com/golang-jwt/jwt/v4"
)

func GetToken(encodedToken string) (string, error) {
	decodedToken, pemKey, err := ValidateToken(encodedToken)

	if err != nil {
		return encodedToken, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, decodedToken)
	log.Printf("new token generated")

	byteKey := []byte(pemKey)
	log.Printf("byteKey generated")
	reEncodedToken, err := token.SignedString(byteKey)

	log.Printf("reEncodedToken: " + reEncodedToken)

	if err != nil {
		return encodedToken, err
	}

	return reEncodedToken, err
}
