package apiKeysRemove_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	apiKeysRemove "github.com/iotea-com/iotea/services/http-api/api/v1/apiKeys/remove"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/api-keys"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		apiKeysRemove.Handler,
		"successfully deletes org API key",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectOrgApiKey(mocks.DB)

			// define prisma mock server response
			testApiKeyId, _ := id.Generator.NewApiKeyId()
			testOrganizationId, _ := id.Generator.NewOrganizationId()

			expectedDatabaseOutput := db.APIKeyModel{}
			mocks.DB.Server.APIKey.Expect(
				mocks.DB.Client.APIKey.FindUnique(
					db.APIKey.ID.Equals(*testApiKeyId),
				),
			).Returns(expectedDatabaseOutput)

			expectedDatabaseOutput = db.APIKeyModel{}
			mocks.DB.Server.APIKey.Expect(
				mocks.DB.Client.APIKey.FindUnique(
					db.APIKey.ID.Equals(*testApiKeyId),
				).Delete(),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"apiKeyId": *testApiKeyId,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithQueryParams := fmt.Sprintf("%s?orgId=%s", route, *testOrganizationId)
			httptestRequest := httptest.NewRequest("DELETE", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			// assert
		},
	)
}
