package services

import (
	"encoding/json"
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jayfaust3/auth.api/pkg/clients"
)

func GetToken(encodedToken string) (string, error) {
	decodedToken, pemKey, err := ValidateToken(encodedToken)

	if err != nil {
		return encodedToken, err
	}

	email := decodedToken.Email

	log.Print("email parsed: " + email)

	user, err := clients.GetUserFromEmail(email)

	if err != nil {
		return encodedToken, err
	} else {
		// userId := user.Id
		// log.Print("User found, id: " + userId)
		encodedUser, encodingError := json.Marshal(user)
		if encodingError != nil {
			log.Print("error occurred while decoding user, error: " + encodingError.Error())
		} else {
			log.Print("user: " + string(encodedUser))
		}

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, decodedToken)

	byteKey := []byte(pemKey)
	reEncodedToken, err := token.SignedString(byteKey)

	if err != nil {
		return encodedToken, err
	}

	return reEncodedToken, err
}
