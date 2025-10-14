package policiesGet_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	policiesGet "github.com/iotea-com/iotea/services/http-api/api/v1/policies/get"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/policies/:certificateId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		policiesGet.Handler,
		"successfully returns certificate",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testCertificateId, _ := id.Generator.NewPermissionSetId()
			expectedDatabaseOutput := db.CertificateModel{
				InnerCertificate: db.InnerCertificate{
					ID:     *testCertificateId,
					Policy: types.JSON{'{', '}'},
				},
			}

			testSpaceId, _ := id.Generator.NewSpaceId()
			mocks.DB.Server.Certificate.Expect(
				mocks.DB.Client.Certificate.FindUnique(
					db.Certificate.ID.Equals(*testCertificateId),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called
			routeWithCertId := strings.ReplaceAll(route, ":certificateId", *testCertificateId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithCertId, *testSpaceId)
			httptestRequest := httptest.NewRequest("GET", routeWithQueryParams, nil)
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// parse response body
			responseBodyBytes, err := io.ReadAll(httptestResponse.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			var response ioteahttp.IoteaApiResponse
			err = json.Unmarshal(responseBodyBytes, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// assert response status and body
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		},
	)
}
