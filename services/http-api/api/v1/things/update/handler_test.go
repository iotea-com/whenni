package thingsUpdate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/id"
	gruentchannel "github.com/ongruent/gruent/libs/legacy/engine/channels"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	thingsUpdate "github.com/ongruent/gruent/services/http-api/api/v1/things/update"
	"github.com/ongruent/gruent/services/http-api/util"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/things/:thingId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		thingsUpdate.Handler,
		"successfully updates thing",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testThingId, _ := id.Generator.NewThingId()
			testThingName := "Test Thing"
			testThingAttributes := map[string]any{
				"test": "attributes",
			}
			testThingAttributesJson, _ := json.Marshal(testThingAttributes)

			expectedDatabaseOutput := db.ThingModel{
				InnerThing: db.InnerThing{
					ID:         *testThingId,
					Name:       testThingName,
					Attributes: testThingAttributesJson, // empty attributes - required by mock Prisma client
				},
			}

			expectedDatabaseChannelOutput := []db.ChannelModel{
				{
					InnerChannel: db.InnerChannel{
						Config: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
					},
				},
			}

			expectedDatabaseThingsOutput := []db.ThingModel{
				{
					InnerThing: db.InnerThing{
						Attributes: types.JSON{'{', '}'}, // empty attributes - required by mock Prisma client
					},
				},
			}

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindMany(
					db.Channel.SpaceID.Equals(*testSpaceId),
				).Select(
					db.Channel.PublishedAt.Field(),
					db.Channel.Config.Field(),
				),
			).ReturnsMany(expectedDatabaseChannelOutput)

			mocks.DB.Server.Thing.Expect(
				mocks.DB.Client.Thing.FindMany(
					db.Thing.SpaceID.Equals(*testSpaceId),
				).Select(
					db.Thing.Attributes.Field(),
				),
			).ReturnsMany(expectedDatabaseThingsOutput)

			mocks.DB.Server.Thing.Expect(
				mocks.DB.Client.Thing.FindUnique(
					db.Thing.ID.Equals(*testThingId),
				).Update(
					db.Thing.Name.Set(testThingName),
					db.Thing.Attributes.Set(testThingAttributesJson),

					// TODO: add user ID to input and set updatedBy
					db.Thing.UpdatedBy.Set("tea_fake"),
					db.Thing.UpdatedAt.Set(util.GetCurrentTime()),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called
			httptestRequestBody := map[string]any{
				"name":       testThingName,
				"attributes": testThingAttributes,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithThingId := strings.ReplaceAll(route, ":thingId", *testThingId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithThingId, *testSpaceId)
			httptestRequest := httptest.NewRequest("PUT", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
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

			// check that Data is a db.ThingModel (marshal into JSON and unmarshal as a db.ThingModel)
			responseBodyJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var thing db.ThingModel
			err = json.Unmarshal(responseBodyJson, &thing)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.ThingModel: %v", err)
			}

			// assert response status and body
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, thing.ID, *testThingId, "thing has the same thing ID as the database output")
		},
	)

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		thingsUpdate.Handler,
		"prevents thing update when in use by a channel",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testThingId, _ := id.Generator.NewThingId()

			exampleConfigWithThing := gruentchannel.Channel{
				Nodes: []gruentchannel.Node{
					{
						Metadata: gruentchannel.NodeMetadata{
							Dependencies: gruentchannel.NodeDependencies{
								Things: []gruentchannel.NodeThingDependency{
									{
										ThingId: *testThingId,
									},
								},
							},
						},
					},
				},
			}

			exampleConfigWithThingJson, _ := json.Marshal(exampleConfigWithThing)

			publishedAt := time.Now()
			expectedDatabaseChannelOutput := []db.ChannelModel{
				{
					InnerChannel: db.InnerChannel{
						PublishedAt: &publishedAt,
						Config:      exampleConfigWithThingJson, // empty config - required by mock Prisma client
					},
				},
			}

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindMany(
					db.Channel.SpaceID.Equals(*testSpaceId),
				).Select(
					db.Channel.PublishedAt.Field(),
					db.Channel.Config.Field(),
				),
			).ReturnsMany(expectedDatabaseChannelOutput)

			// run integration test
			httptestRequestBody := map[string]any{
				"name": "Test Thing",
				"attributes": map[string]any{
					"test": "attributes",
				},
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithThingId := strings.ReplaceAll(route, ":thingId", *testThingId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithThingId, *testSpaceId)
			httptestRequest := httptest.NewRequest("DELETE", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, err := mocks.API.App.Test(httptestRequest)

			// assert nil error
			assert.Nil(t, err)

			// assert response status
			assert.Equalf(t, http.StatusConflict, httptestResponse.StatusCode, "matching response code")
		},
	)
}
