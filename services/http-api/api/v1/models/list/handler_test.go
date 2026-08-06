package modelsList_test

import (
	// "encoding/json"
	// "io"
	// "net/http"
	// "net/http/httptest"
	// "strings"
	"testing"
	// gruenthttp "github.com/ongruent/gruent/libs/http"
	// "github.com/ongruent/gruent/libs/id"
	// "github.com/ongruent/gruent/prisma/db"
	// apitest "github.com/ongruent/gruent/services/http-api/api/test"
	// structsList "github.com/ongruent/gruent/services/http-api/api/v1/spaces/structs/list"
	// "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// route := "/v1/spaces/:spaceId/structs"

	// TODO: re-enable list test when Prisma is replaced
	// apitest.HandlerUnitTestWithSetup(
	// 	t,
	// 	route,
	// 	structsList.Handler,
	// 	"successfully returns list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		testStructId, _ := id.Generator.NewStructId()
	// 		expectedDatabaseOutput := []db.StructModel{
	// 			{
	// 				InnerStruct: db.InnerStruct{
	// 					ID:     *testStructId,
	// 					Schema: types.JSON{'{', '}'}, // empty schema - required by mock Prisma client
	// 				},
	// 			},
	// 		}

	// 		mocks.DB.Server.Struct.Expect(
	// 			mocks.DB.Client.Struct.FindMany(
	// 				db.Struct.SpaceID.Equals(*testSpaceId),
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

	// 		var response gruenthttp.GruentApiResponse
	// 		err = json.Unmarshal(responseBodyBytes, &response)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal response body: %v", err)
	// 		}

	// 		// check that Data is a []db.StructModel (marshal into JSON and unmarshal as a []db.StructModel)
	// 		responseBodyJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var structs []db.StructModel
	// 		err = json.Unmarshal(responseBodyJson, &structs)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into []db.StructModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, structs, 1, "structs list has a length of 1")
	// 		assert.Equal(t, structs[0].ID, *testStructId, "structs[0] has the same space ID as the database output")
	// 	},
	// )

	// apitest.HandlerUnitTestWithSetup(t,
	// 	route,
	// 	structsList.Handler,
	// 	"successfully returns empty list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks,
	// 	) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		mocks.DB.Server.Struct.Expect(
	// 			mocks.DB.Client.Struct.FindMany(
	// 				db.Struct.SpaceID.Equals(*testSpaceId),
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

	// 		var response gruenthttp.GruentApiResponse
	// 		err = json.Unmarshal(responseBodyBytes, &response)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal response body: %v", err)
	// 		}

	// 		// check that Data is a []db.StructModel (marshal into JSON and unmarshal as a []db.StructModel)
	// 		responseBodyJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var structs []db.StructModel
	// 		err = json.Unmarshal(responseBodyJson, &structs)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into []db.StructModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, structs, 0, "length of structs list is 0")
	// 	})
}
