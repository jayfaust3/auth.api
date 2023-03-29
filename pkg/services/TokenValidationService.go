package services

import (
	"github.com/golang-jwt/jwt"
	"../models/application/token/TokenPayload.go"
	"./SigningCertFetchService.go"
)

func validateToken(encodedToken: string) (TokenPayload, error) {

}
