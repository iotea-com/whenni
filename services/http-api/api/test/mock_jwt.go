package apitest

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/iotea-com/iotea/services/http-api/config"
)

func NewMockJwt(t *testing.T, testUserId string) string {
	testUserJwt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": testUserId,
	}).SignedString([]byte(config.VaultConf.JwtSecret))
	if err != nil {
		t.Fatalf("could not sign jwt: %s", err)
	}

	return testUserJwt
}
