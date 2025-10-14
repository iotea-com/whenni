package apitest

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func AssertRouteExists(t *testing.T, routes []fiber.Route, method string, path string) {
	exists := false
	for _, r := range routes {
		if r.Method == method && r.Path == path {
			exists = true
			break
		}
	}
	assert.True(t, exists, "route with method %s and path %s should be registered", method, path)
}
