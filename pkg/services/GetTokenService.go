package services

import "fmt"

func GetToken(encodedToken string) (string, error) {
	decodedToken, err := ValidateToken(encodedToken)

	fmt.Print("decodedToken: ", decodedToken)

	return encodedToken, err
}
