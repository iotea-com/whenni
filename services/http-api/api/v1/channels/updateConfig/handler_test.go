package channelsUpdateConfig_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	channelsUpdateConfig "github.com/ongruent/gruent/services/http-api/api/v1/channels/updateConfig"
	"github.com/ongruent/gruent/services/http-api/util"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/channels/:channelId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		channelsUpdateConfig.Handler,
		"successfully updates channel config",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testChannelId, _ := id.Generator.NewChannelId()
			testChannelConfig := types.JSON{'{', '}'}
			expectedDatabaseOutput := db.ChannelModel{
				InnerChannel: db.InnerChannel{
					ID:     *testChannelId,
					Config: testChannelConfig,
				},
			}

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindUnique(
					db.Channel.ID.Equals(*testChannelId),
				).Update(
					db.Channel.Config.Set(testChannelConfig),
					db.Channel.UpdatedAt.Set(util.GetCurrentTime()),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			routeWithChannelId := strings.ReplaceAll(route, ":channelId", *testChannelId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithChannelId, *testSpaceId)
			httptestRequest := httptest.NewRequest("PUT", routeWithQueryParams, bytes.NewBuffer(testChannelConfig))
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// parse response body
			responseBodyBytes, err := io.ReadAll(httptestResponse.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			var response gruenthttp.GruentApiResponse
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

			// assert response status and body
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, channel.ID, *testChannelId, "channel has the same channel ID as the database output")
		},
	)
}
