package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jayfaust3/auth.api/pkg/models/application/token"
)

func validateToken(encodedToken string) (token.TokenPayload, error) {
	claimsStruct := token.TokenPayload{}

	token, err := jwt.ParseWithClaims(
		encodedToken,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			keyId := fmt.Sprintf("%s", token.Header["kid"])

			pem, err := getSigningCert(keyId)

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
		return token.TokenPayload{}, err
	}

	claims, ok := token.Claims.(*token.TokenPayload)

	if !ok {
		return token.TokenPayload{}, errors.New("Invalid token")
	}

	if claims.Issuer != os.Getenv("AUTH_ISSUER") {
		return token.TokenPayload{}, errors.New("iss is invalid")
	}

	if claims.Audience != os.Getenv("AUTH_AUDIENCE") {
		return token.TokenPayload{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return token.TokenPayload{}, errors.New("JWT is expired")
	}

	return *claims, nil
}
