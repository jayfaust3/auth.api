package token

import (
	"jwt"
)

type TokenPayload struct {
	UserId        string `json:"userId"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.StandardClaims
}
