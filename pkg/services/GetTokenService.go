package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jayfaust3/auth.api/pkg/clients"
)

func GetToken(encodedToken string) (string, error) {
	decodedToken, pemKey, err := ValidateToken(encodedToken)

	if err != nil {
		return encodedToken, err
	}

	decodedToken.Id = uuid.New().String()
	email := decodedToken.Email

	log.Print("email parsed: " + email)

	user, err := clients.GetUserFromEmail(email)

	if err == nil {
		userId := user.Id
		log.Print(fmt.Sprintf("User found, id: %s", userId))

		permissions, err := clients.GetPermissionsByEntity(userId, 0)

		if err == nil {
			scopes := []string{}

			for _, permission := range permissions {
				resource := permission.Resource
				action := permission.Action
				scopes = append(scopes, fmt.Sprintf("%s:%s", resource, action))
			}

			scope := strings.Join(scopes, " ")

			log.Print(fmt.Sprintf("Applying the following scopes to the token: %s", scope))

			decodedToken.Scope = scope
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, decodedToken)

	byteKey := []byte(pemKey)
	reEncodedToken, err := token.SignedString(byteKey)

	return reEncodedToken, err
}
