package ioteahttputil

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJwt(secret string) func(token *jwt.Token) (any, error) {
	return func(token *jwt.Token) (any, error) {
		// validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// return secret
		return []byte(secret), nil
	}
}
