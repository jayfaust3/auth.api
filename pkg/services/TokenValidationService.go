package services

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jayfaust3/auth.api/pkg/models/application/token"
)

type TokenPayload = token.TokenPayload

func ValidateToken(encodedToken string) (token.TokenPayload, rsa.PublicKey, error) {
	claimsStruct := token.TokenPayload{}

	var publicKey rsa.PublicKey

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

			publicKey = *key

			if err != nil {
				return nil, err
			}

			return key, nil
		},
	)

	if err != nil {
		return TokenPayload{}, publicKey, err
	}

	claims, ok := token.Claims.(*TokenPayload)

	if !ok {
		return TokenPayload{}, publicKey, errors.New("Invalid token")
	}

	if claims.Issuer != os.Getenv("AUTH_ISSUER") {
		return TokenPayload{}, publicKey, errors.New("iss is invalid")
	}

	if claims.Audience != os.Getenv("AUTH_AUDIENCE") {
		return TokenPayload{}, publicKey, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return TokenPayload{}, publicKey, errors.New("JWT is expired")
	}

	return *claims, publicKey, nil
}
