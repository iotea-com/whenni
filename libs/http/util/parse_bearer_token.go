package ioteahttputil

import (
	"fmt"
	"strings"
)

func ParseBearerToken(authHeader string) (*string, error) {
	splitAuthHeader := strings.Split(authHeader, " ")

	if len(splitAuthHeader) < 2 {
		err := fmt.Errorf("the provided string was not a valid auth header with a bearer token value")
		return nil, err
	}

	if splitAuthHeader[0] != "Bearer" {
		err := fmt.Errorf("the provided string was not a valid auth header with a bearer token value")
		return nil, err
	}

	token := splitAuthHeader[1]

	return &token, nil
}
