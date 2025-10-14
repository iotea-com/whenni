package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	t.Run("successfully creates a new api", func(t *testing.T) {
		// when
		// a new API is created
		api := New()

		// then
		// assert HTTP server exists and there are routes available
		assert.NotNil(t, api, "API exists")
		assert.NotNil(t, api.HttpServer, "HTTP server exists on the API")
		assert.Greater(t, len(api.HttpServer.GetRoutes()), 0, "HTTP server has routes available")
	})
}
