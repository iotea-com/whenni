package get_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ongruent/gruent/services/http-api/api/healthcheck/get"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/healthcheck"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		get.Handler,
		"successfully returns 200",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// when
			// API route is called
			httptestRequest := httptest.NewRequest("GET", route, nil)

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		},
	)
}
