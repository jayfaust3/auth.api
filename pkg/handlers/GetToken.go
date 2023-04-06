package handlers

import (
	"net/http"
	"strings"

	"github.com/jayfaust3/auth.api/pkg/api/responses"
	"github.com/jayfaust3/auth.api/pkg/api/responses/auth"
	"github.com/jayfaust3/auth.api/pkg/services"
	"github.com/jayfaust3/auth.api/pkg/utils"
)

func GetToken(w http.ResponseWriter, r *http.Request) {
	authHeaderValue := r.Header.Get("Authorization")

	var encodedToken string

	if strings.Contains(authHeaderValue, "Bearer") {
		encodedToken = strings.ReplaceAll(authHeaderValue, "Bearer ", "")
	} else {
		encodedToken = authHeaderValue
	}

	generatedToken, err := services.GetToken(encodedToken)

	if err != nil {
		var apiResponse responses.ApiResponse[auth.AuthTokenResponse]

		var tokenResponse auth.AuthTokenResponse

		tokenResponse.AuthToken = generatedToken

		apiResponse.Data = tokenResponse

		utils.RespondWithJSON(w, 200, apiResponse)
	} else {
		utils.RespondWithError(w, 400, "Unable to verify token")
	}
}
