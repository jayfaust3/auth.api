package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jayfaust3/auth.api/pkg/models/application/token"
)

type TokenPayload = token.TokenPayload

func ValidateToken(encodedToken string) (token.TokenPayload, string, error) {
	claimsStruct := token.TokenPayload{}

	var pemKey string

	token, err := jwt.ParseWithClaims(
		encodedToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			keyId := fmt.Sprintf("%s", token.Header["kid"])

			pem, err := GetSigningCert(keyId)
			pemKey = pem

			if err != nil {
				return nil, err
			}

			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))

			if err != nil {
				return nil, err
			}

			return key, nil
		},
	)

	if err != nil {
		return TokenPayload{}, pemKey, err
	}

	claims, ok := token.Claims.(*TokenPayload)

	if !ok {
		return TokenPayload{}, pemKey, errors.New("Invalid token")
	}

	if claims.Issuer != os.Getenv("AUTH_ISSUER") {
		return TokenPayload{}, pemKey, errors.New("iss is invalid")
	}

	if claims.Audience != os.Getenv("AUTH_AUDIENCE") {
		return TokenPayload{}, pemKey, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return TokenPayload{}, pemKey, errors.New("JWT is expired")
	}

	return *claims, pemKey, nil
}
