package apiKeysAdd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	apiKeysAdd "github.com/iotea-com/iotea/services/http-api/api/v1/apiKeys/add"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/api-keys"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		apiKeysAdd.Handler,
		"successfully inserts org API key",
		func(t *testing.T, mocks *apitest.HandlerMocks,
		) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectOrgApiKey(mocks.DB)

			// define prisma mock server response
			testApiKeyId, _ := id.Generator.NewApiKeyId()
			testApiKeyName, _ := id.Generator.NewApiKeyId()
			testOrganizationId, _ := id.Generator.NewOrganizationId()
			testPermissionSetId, _ := id.Generator.NewPermissionSetId()

			expectedDatabaseOutput := db.APIKeyModel{
				InnerAPIKey: db.InnerAPIKey{
					ID:                          *testApiKeyId,
					OrganizationID:              *testOrganizationId,
					OrganizationPermissionSetID: testPermissionSetId,
					CreatedBy:                   "tea_fake",
					CreatedAt:                   time.Now(),
					ExpiresAt:                   nil,
				},
			}

			mocks.DB.Server.APIKey.Expect(
				mocks.DB.Client.APIKey.CreateOne(
					db.APIKey.ID.Set(*testApiKeyId),
					db.APIKey.CreatedBy.Set("tea_fake"),
					db.APIKey.Organization.Link(
						db.Organization.ID.Equals(*testOrganizationId),
					),
					db.APIKey.OrganizationPermissionSet.Link(
						db.PermissionSet.ID.Equals(*testPermissionSetId),
					),
					db.APIKey.Name.SetIfPresent(testApiKeyName),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"permissionSetId": *testPermissionSetId,
				"name":            testApiKeyName,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithQueryParams := fmt.Sprintf("%s?orgId=%s", route, *testOrganizationId)
			httptestRequest := httptest.NewRequest("POST", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
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

			// check that Data is a db.APIKeyModel struct (marshal into JSON and unmarshal as an apiKey)
			apiKeyJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var apiKey db.APIKeyModel
			err = json.Unmarshal(apiKeyJson, &apiKey)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.APIKeyModel: %v", err)
			}

			// assert response status and body
			if httptestResponse.StatusCode != http.StatusCreated {
				t.Logf("Expected status code %d, but got %d", http.StatusCreated, httptestResponse.StatusCode)
				t.Fatalf("%v", response.Errors[0])
			}

			assert.Equalf(t, http.StatusCreated, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, *testApiKeyId, apiKey.ID, "API Key ID should match the expected value")
		})

	// Future tests - TODO:
	// case for non-existent permission set
	// case for permission set is not in the space
}
