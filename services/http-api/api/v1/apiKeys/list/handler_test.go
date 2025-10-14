package apiKeysList_test

import (
	// "encoding/json"
	// "io"
	// "net/http"
	// "net/http/httptest"
	// "strings"
	"testing"
	// ioteahttp "github.com/iotea-com/iotea/libs/http"
	// "github.com/iotea-com/iotea/libs/id"
	// "github.com/iotea-com/iotea/prisma/db"
	// apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	// apiKeysList "github.com/iotea-com/iotea/services/http-api/api/v1/spaces/apiKeys/list"
	// "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// route := "/v1/spaces/:spaceId/api-keys"

	// TODO: re-enable list test when Prisma is replaced
	// apitest.HandlerUnitTestWithSetup(
	// 	t,
	// 	route,
	// 	apiKeysList.Handler,
	// 	"successfully returns list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testApiKeyId, _ := id.Generator.NewApiKeyId()
	// 		expectedDatabaseOutput := []db.APIKeyModel{
	// 			{
	// 				InnerAPIKey: db.InnerAPIKey{
	// 					ID: *testApiKeyId,
	// 				},
	// 			},
	// 		}

	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		mocks.DB.Server.APIKey.Expect(
	// 			mocks.DB.Client.APIKey.FindMany(
	// 				db.APIKey.SpaceID.Equals(*testSpaceId),
	// 			).With(
	// 				db.APIKey.PermissionSet.Fetch(),
	// 			),
	// 		).ReturnsMany(expectedDatabaseOutput)

	// 		// when
	// 		// API route is called
	// 		routeWithParams := strings.ReplaceAll(route, ":spaceId", *testSpaceId)
	// 		httptestRequest := httptest.NewRequest("GET", routeWithParams, nil)
	// 		httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

	// 		httptestResponse, _ := mocks.API.App.Test(httptestRequest)

	// 		// then
	// 		// parse response body
	// 		responseBodyBytes, err := io.ReadAll(httptestResponse.Body)
	// 		if err != nil {
	// 			t.Fatalf("Failed to read response body: %v", err)
	// 		}

	// 		var response ioteahttp.IoteaApiResponse
	// 		err = json.Unmarshal(responseBodyBytes, &response)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal response body: %v", err)
	// 		}

	// 		// check that Data is a db.APIKeyModel struct (marshal into JSON and unmarshal as an apiKey)
	// 		apiKeysJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var apiKeys []db.APIKeyModel
	// 		err = json.Unmarshal(apiKeysJson, &apiKeys)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into db.APIKeyModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, apiKeys, 1, "apiKeys list has a length of 1")
	// 		assert.Equal(t, apiKeys[0].ID, *testApiKeyId, "apiKeys[0] has the same ID as the database output")
	// 	},
	// )

	// apitest.HandlerUnitTestWithSetup(t,
	// 	route,
	// 	apiKeysList.Handler,
	// 	"successfully returns empty list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks,
	// 	) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		mocks.DB.Server.APIKey.Expect(
	// 			mocks.DB.Client.APIKey.FindMany(
	// 				db.APIKey.SpaceID.Equals(*testSpaceId),
	// 			).With(
	// 				db.APIKey.PermissionSet.Fetch(),
	// 			),
	// 		).Errors(db.ErrNotFound)

	// 		// when
	// 		// API route is called and no results are found in the database
	// 		routeWithParams := strings.ReplaceAll(route, ":spaceId", *testSpaceId)
	// 		httptestRequest := httptest.NewRequest("GET", routeWithParams, nil)
	// 		httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

	// 		httptestResponse, _ := mocks.API.App.Test(httptestRequest)

	// 		// then
	// 		// parse response body
	// 		responseBodyBytes, err := io.ReadAll(httptestResponse.Body)
	// 		if err != nil {
	// 			t.Fatalf("Failed to read response body: %v", err)
	// 		}

	// 		var response ioteahttp.IoteaApiResponse
	// 		err = json.Unmarshal(responseBodyBytes, &response)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal response body: %v", err)
	// 		}

	// 		// check that Data is a db.APIKeyModel struct (marshal into JSON and unmarshal as an apiKey)
	// 		apiKeysJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var apiKeys []db.APIKeyModel
	// 		err = json.Unmarshal(apiKeysJson, &apiKeys)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into db.APIKeyModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, apiKeys, 0, "apiKeys list has a length of 0")
	// 	})
}
