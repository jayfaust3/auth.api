package handlers

import (
	"encoding/json"
	"net/http"
)

func getToken(w http.ResponseWriter, r *http.Request) {
	tokenFromHeader := r.Header.Get("Authorization")

	generatedToken, err := services.getToken(tokenFromHeader)

	respondWithJSON(w, 200, struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}{
		Data: {
			Token: generatedToken,
		},
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}
