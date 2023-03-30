package services

import (
	"github.com/jayfaust3/auth.api/pkg/models/application/token"
)

func getToken(encodedToken: string) (string, error) {
	decodedToken := validateToken(encodedToken)

	return encodedToken, nil
}
