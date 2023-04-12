package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

func GetSigningCert(keyID string) (string, error) {
	signingCertEndpoint := os.Getenv("GOOGLE_AUTH_CERT_ENDPOINT")

	if signingCertEndpoint == "" {
		return "", errors.New("Google auth cert endpoint not configured")
	}

	resp, err := http.Get(signingCertEndpoint)

	if err != nil {
		return "", err
	}

	dat, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	certDict := map[string]string{}
	err = json.Unmarshal(dat, &certDict)

	if err != nil {
		return "", err
	}

	key, ok := certDict[keyID]

	if !ok {
		return "", errors.New("key not found")
	}

	return key, nil
}
