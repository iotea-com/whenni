package channelsGet_test

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
	channelsGet "github.com/iotea-com/iotea/services/http-api/api/v1/channels/get"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/channels/:channelId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		channelsGet.Handler,
		"successfully gets channel",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testChannelId, _ := id.Generator.NewChannelId()
			expectedDatabaseOutput := db.ChannelModel{
				InnerChannel: db.InnerChannel{
					ID:     *testChannelId,
					Config: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
				},
			}
			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindUnique(
					db.Channel.ID.Equals(*testChannelId),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			routeWithChannelId := strings.ReplaceAll(route, ":channelId", *testChannelId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithChannelId, *testSpaceId)
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

			// check that Data is a db.ChannelModel (marshal into JSON and unmarshal as a db.ChannelModel)
			responseBodyJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var channel db.ChannelModel
			err = json.Unmarshal(responseBodyJson, &channel)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.ChannelModel: %v", err)
			}

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Equalf(t, channel.ID, *testChannelId, "channel ID matches database output")
		},
	)
}
