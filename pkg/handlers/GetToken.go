package handlers

import (
	"encoding/json"
	"github.com/jayfaust3/auth.api/pkg/services"
	"net/http"
	"strconv"
)

func getToken(w http.ResponseWriter, r *http.Request) {
	tokenFromHeader := r.Header.Get("Authorization")

	generatedToken, err := services.getToken(tokenFromHeader)

	respondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{
		Token: generatedToken,
	})
}
