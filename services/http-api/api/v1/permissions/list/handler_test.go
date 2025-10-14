package permissionsList_test

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
	// "github.com/iotea-com/iotea/services/http-api/api/v1/spaces/permissions/list"
	// "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// route := "/v1/spaces/:spaceId/permissions"

	// TODO: re-enable list test when Prisma is replaced
	// apitest.HandlerUnitTestWithSetup(
	// 	t,
	// 	route,
	// 	permissionsList.Handler,
	// 	"successfully returns list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testPermissionSetId, _ := id.Generator.NewPermissionSetId()
	// 		expectedDatabaseOutput := []db.PermissionSetModel{
	// 			{
	// 				InnerPermissionSet: db.InnerPermissionSet{
	// 					ID: *testPermissionSetId,
	// 				},
	// 			},
	// 		}

	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		mocks.DB.Server.PermissionSet.Expect(
	// 			mocks.DB.Client.PermissionSet.FindMany(
	// 				db.PermissionSet.SpaceID.Equals(*testSpaceId),
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

	// 		// check that Data is a db.PermissionSetModel struct (marshal into JSON and unmarshal as a permission set)
	// 		permissionSetsJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var permissionSets []db.PermissionSetModel
	// 		err = json.Unmarshal(permissionSetsJson, &permissionSets)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into db.PermissionSetModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, permissionSets, 1, "permissionSets list has a length of 1")
	// 		assert.Equal(t, permissionSets[0].ID, *testPermissionSetId, "permissionSets[0] has the same ID as the database output")
	// 	},
	// )

	// apitest.HandlerUnitTestWithSetup(t,
	// 	route,
	// 	permissionsList.Handler,
	// 	"successfully returns empty list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks,
	// 	) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		mocks.DB.Server.PermissionSet.Expect(
	// 			mocks.DB.Client.PermissionSet.FindMany(
	// 				db.PermissionSet.SpaceID.Equals(*testSpaceId),
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

	// 		// check that Data is a db.PermissionSetModel struct (marshal into JSON and unmarshal as a permission set)
	// 		permissionSetsJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var permissionSets []db.PermissionSetModel
	// 		err = json.Unmarshal(permissionSetsJson, &permissionSets)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into db.PermissionSetModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, permissionSets, 0, "permissionSets list has a length of 0")
	// 	})
}
