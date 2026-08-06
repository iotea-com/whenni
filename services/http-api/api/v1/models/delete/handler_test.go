package modelsDelete_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ongruent/gruent/libs/id"
	gruentchannels "github.com/ongruent/gruent/libs/legacy/engine/channels"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	modelsDelete "github.com/ongruent/gruent/services/http-api/api/v1/models/delete"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/models/:modelId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		modelsDelete.Handler,
		"successfully deletes model",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testModelId, _ := id.Generator.NewModelId()

			expectedDatabaseChannelOutput := []db.ChannelModel{
				{
					InnerChannel: db.InnerChannel{
						Config: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
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

			expectedDatabaseModelOutput := db.ModelModel{
				InnerModel: db.InnerModel{
					Attributes: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
				},
			}

			mocks.DB.Server.Model.Expect(
				mocks.DB.Client.Model.FindUnique(
					db.Model.ID.Equals(*testModelId),
				).Delete(),
			).Returns(expectedDatabaseModelOutput)

			// run integration test
			routeWithModelId := strings.ReplaceAll(route, ":modelId", *testModelId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithModelId, *testSpaceId)
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
		modelsDelete.Handler,
		"prevents model deletion when in use by a channel",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testModelId, _ := id.Generator.NewModelId()

			exampleConfigWithStruct := gruentchannels.Channel{
				Nodes: []gruentchannels.Node{
					{
						Metadata: gruentchannels.NodeMetadata{
							Dependencies: gruentchannels.NodeDependencies{
								Models: []gruentchannels.NodeModelDependency{
									{
										ModelId: *testModelId,
									},
								},
							},
						},
					},
				},
			}

			exampleConfigWithStructJson, _ := json.Marshal(exampleConfigWithStruct)

			expectedDatabaseChannelOutput := []db.ChannelModel{
				{
					InnerChannel: db.InnerChannel{
						Config: exampleConfigWithStructJson, // empty config - required by mock Prisma client
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
			routeWithModelId := strings.ReplaceAll(route, ":modelId", *testModelId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithModelId, *testSpaceId)
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
