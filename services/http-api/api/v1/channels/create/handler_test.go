package channelsCreate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	channelsCreate "github.com/ongruent/gruent/services/http-api/api/v1/channels/create"
	"github.com/ongruent/gruent/services/http-api/util"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/channels"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		channelsCreate.Handler,
		"successfully creates channel",
		func(t *testing.T, mocks *apitest.HandlerMocks,
		) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testChannelName := "Test Channel"
			testChannelId, _ := id.Generator.NewChannelId()
			testChannelConfig := map[string]any{
				"id": *testChannelId,
			}
			testThingId := testChannelId
			testThingName := testChannelName

			marshalledTestChannelConfig, _ := json.Marshal(testChannelConfig)

			expectedDatabaseChannelOutput := db.ChannelModel{
				InnerChannel: db.InnerChannel{
					ID:     *testChannelId,
					Config: marshalledTestChannelConfig,
				},
			}

			expectedDatabaseThingOutput := db.ThingModel{
				InnerThing: db.InnerThing{
					ID:         *testThingId,
					Attributes: db.JSON{'{', '}'},
				},
			}

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.CreateOne(
					db.Channel.ID.Set(*testChannelId),
					db.Channel.Name.Set(testChannelName),
					db.Channel.Config.Set(marshalledTestChannelConfig),
					db.Channel.CreatedBy.Set("tea_fake"),
					db.Channel.UpdatedBy.Set("tea_fake"),
					db.Channel.Space.Link(
						db.Space.ID.Equals(*testSpaceId),
					),
				),
			).Returns(expectedDatabaseChannelOutput)

			now := util.GetCurrentTime()
			mocks.DB.Server.Thing.Expect(
				mocks.DB.Client.Thing.CreateOne(
					db.Thing.ID.Set(*testThingId),
					db.Thing.Name.Set(testThingName),
					db.Thing.Attributes.Set(db.JSON{'{', '}'}),
					db.Thing.CreatedBy.Set("internal"),
					db.Thing.UpdatedBy.Set("internal"),
					db.Thing.ThingCategory.Set(things.HttpServerThingCategory.String()),
					db.Thing.Space.Link(
						db.Space.ID.Equals(*testSpaceId),
					),
					db.Thing.Internal.Set(true),
					db.Thing.CreatedAt.Set(now),
					db.Thing.UpdatedAt.Set(now),
				),
			).Returns(expectedDatabaseThingOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"name":   testChannelName,
				"config": testChannelConfig,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", mocks.API.Route, *testSpaceId)
			httptestRequest := httptest.NewRequest("POST", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
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
			assert.Equalf(t, http.StatusCreated, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, channel.ID, *testChannelId, "channel has the same ID as the database output")
		})

	// Future tests - TODO:
	// case for channel limit in a space
}
