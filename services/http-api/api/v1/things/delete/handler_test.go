package thingsDelete_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ioteachannel "github.com/iotea-com/iotea/libs/engine/channels"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	thingsDelete "github.com/iotea-com/iotea/services/http-api/api/v1/things/delete"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/things/:thingId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		thingsDelete.Handler,
		"successfully deletes thing",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testThingId, _ := id.Generator.NewThingId()

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
						Attributes: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
					},
				},
			}

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindMany(
					db.Channel.SpaceID.Equals(*testSpaceId),
				).Select(
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

			expectedDatabaseThingOutput := db.ThingModel{
				InnerThing: db.InnerThing{
					Attributes: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
				},
			}

			mocks.DB.Server.Thing.Expect(
				mocks.DB.Client.Thing.FindUnique(
					db.Thing.ID.Equals(*testThingId),
				).Delete(),
			).Returns(expectedDatabaseThingOutput)

			// run integration test
			routeWithThingId := strings.ReplaceAll(route, ":thingId", *testThingId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithThingId, *testSpaceId)
			httptestRequest := httptest.NewRequest("DELETE", routeWithQueryParams, nil)
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, err := mocks.API.App.Test(httptestRequest)

			// assert nil error
			assert.Nil(t, err)

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		},
	)

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		thingsDelete.Handler,
		"prevents thing deletion when in use by a channel",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testThingId, _ := id.Generator.NewThingId()

			exampleConfigWithThing := ioteachannel.Channel{
				Nodes: []ioteachannel.Node{
					{
						Metadata: ioteachannel.NodeMetadata{
							Dependencies: ioteachannel.NodeDependencies{
								Things: []ioteachannel.NodeThingDependency{
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

			expectedDatabaseChannelOutput := []db.ChannelModel{
				{
					InnerChannel: db.InnerChannel{
						Config: exampleConfigWithThingJson, // empty config - required by mock Prisma client
					},
				},
			}

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindMany(
					db.Channel.SpaceID.Equals(*testSpaceId),
				).Select(
					db.Channel.Config.Field(),
				),
			).ReturnsMany(expectedDatabaseChannelOutput)

			// run integration test
			routeWithThingId := strings.ReplaceAll(route, ":thingId", *testThingId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithThingId, *testSpaceId)
			httptestRequest := httptest.NewRequest("DELETE", routeWithQueryParams, nil)
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, err := mocks.API.App.Test(httptestRequest)

			// assert nil error
			assert.Nil(t, err)

			// assert response status
			assert.Equalf(t, http.StatusConflict, httptestResponse.StatusCode, "matching response code")
		},
	)
}
