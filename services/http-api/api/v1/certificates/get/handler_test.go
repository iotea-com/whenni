package certificatesGet_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ongruent/gruent/libs/id"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	certificatesGet "github.com/ongruent/gruent/services/http-api/api/v1/certificates/get"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/policies/:certificateId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		certificatesGet.Handler,
		"successfully returns certificate",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define secrets mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()

			testCertificateId := "683902d1126b74151b16aec61b70a94bf662f7fa"
			expectedCertificatesKeyOutput := "fake-key"
			expectedCertificatesCertOutput := "fake-cert"
			mocks.Secrets.Client.Mock.On("RetrieveCert", testCertificateId).
				Return(&expectedCertificatesCertOutput, &expectedCertificatesKeyOutput, nil)

			expectedRootCaCertificateCertOutput := "fake-cert"
			mocks.Secrets.Client.Mock.On("RetrieveRootCaCert").
				Return(&expectedRootCaCertificateCertOutput, nil)

			// when
			// API route is called
			routeWithParams := strings.ReplaceAll(route, ":certificateId", testCertificateId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithParams, *testSpaceId)
			httptestRequest := httptest.NewRequest("GET", routeWithQueryParams, nil)
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		},
	)
}
