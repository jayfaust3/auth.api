package token

import (
	"github.com/golang-jwt/jwt/v4"
)

type TokenPayload struct {
	UserId        string `json:"userId"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Scope         string `json:"scp"`
	jwt.StandardClaims
}
