package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func GetSigningCert(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
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
