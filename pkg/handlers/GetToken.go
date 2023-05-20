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

	encodedToken := strings.ReplaceAll(authHeaderValue, "Bearer ", "")

	generatedToken, err := services.GetToken(encodedToken)

	if err == nil {
		apiResponse := responses.ApiResponse[auth.AuthTokenResponse]{
			Data: auth.AuthTokenResponse{
				AuthToken: generatedToken,
			},
		}

		utils.RespondWithJSON(w, 200, apiResponse)
	} else {
		utils.RespondWithError(w, 400, "Unable to verify token, error: "+err.Error())
	}
}
