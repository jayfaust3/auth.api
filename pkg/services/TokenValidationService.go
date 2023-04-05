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

func ValidateToken(encodedToken string) (TokenPayload, error) {
	claimsStruct := token.TokenPayload{}

	token, err := jwt.ParseWithClaims(
		encodedToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			keyId := fmt.Sprintf("%s", token.Header["kid"])

			pem, err := GetSigningCert(keyId)

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
		return TokenPayload{}, err
	}

	claims, ok := token.Claims.(*TokenPayload)

	if !ok {
		return TokenPayload{}, errors.New("Invalid token")
	}

	if claims.Issuer != os.Getenv("AUTH_ISSUER") {
		return TokenPayload{}, errors.New("iss is invalid")
	}

	if claims.Audience != os.Getenv("AUTH_AUDIENCE") {
		return TokenPayload{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return TokenPayload{}, errors.New("JWT is expired")
	}

	return *claims, nil
}
