package channelsList_test

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
	// channelsList "github.com/ongruent/gruent/services/http-api/api/v1/spaces/channels/list"
	// "github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	// route := "/v1/spaces/:spaceId/channels"

	// TODO: re-enable list test when Prisma is replaced
	// apitest.HandlerUnitTestWithSetup(
	// 	t,
	// 	route,
	// 	channelsList.Handler,
	// 	"successfully returns list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		testChannelId, _ := id.Generator.NewChannelId()
	// 		expectedDatabaseOutput := []db.ChannelModel{
	// 			{
	// 				InnerChannel: db.InnerChannel{
	// 					ID:     *testChannelId,
	// 					Config: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
	// 				},
	// 			},
	// 		}

	// 		mocks.DB.Server.Channel.Expect(
	// 			mocks.DB.Client.Channel.FindMany(
	// 				db.Channel.SpaceID.Equals(*testSpaceId),
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

	// 		// check that Data is a []db.ChannelModel struct (marshal into JSON and unmarshal as a []db.Channel)
	// 		channelsJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var channels []db.ChannelModel
	// 		err = json.Unmarshal(channelsJson, &channels)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into db.PermissionSetModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, channels, 1, "channels list has a length of 1")
	// 		assert.Equal(t, channels[0].ID, *testChannelId, "channels[0] has the same space ID as the database output")
	// 	},
	// )

	// apitest.HandlerUnitTestWithSetup(t,
	// 	route,
	// 	channelsList.Handler,
	// 	"successfully returns empty list",
	// 	func(t *testing.T, mocks *apitest.HandlerMocks,
	// 	) {
	// 		// given
	// 		// define prisma mock server response for bearer token
	// 		apitest.SetupExpectApiKey(mocks.DB)

	// 		// define prisma mock server response
	// 		testSpaceId, _ := id.Generator.NewSpaceId()
	// 		mocks.DB.Server.Channel.Expect(
	// 			mocks.DB.Client.Channel.FindMany(
	// 				db.Channel.SpaceID.Equals(*testSpaceId),
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

	// 		// check that Data is a []db.ChannelModel struct (marshal into JSON and unmarshal as a []db.Channel)
	// 		channelsJson, err := json.Marshal(response.Data)
	// 		if err != nil {
	// 			t.Fatalf("Failed to marshal response data back to JSON: %v", err)
	// 		}

	// 		var channels []db.ChannelModel
	// 		err = json.Unmarshal(channelsJson, &channels)
	// 		if err != nil {
	// 			t.Fatalf("Failed to unmarshal data into db.PermissionSetModel: %v", err)
	// 		}

	// 		// assert response status and body
	// 		assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
	// 		assert.Len(t, channels, 0, "length of channels list is 0")
	// 	})
}
