package services

import (
	"log"

	"github.com/golang-jwt/jwt/v4"
)

func GetToken(encodedToken string) (string, error) {
	decodedToken, publicKey, err := ValidateToken(encodedToken)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, decodedToken)
	log.Printf("new token generated")

	reEncodedToken, err := token.SignedString(publicKey)
	log.Printf("error occurred, error: " + err.Error())

	return reEncodedToken, err
}
