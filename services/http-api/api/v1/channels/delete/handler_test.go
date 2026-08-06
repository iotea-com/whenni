package channelsDelete_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	channelsDelete "github.com/ongruent/gruent/services/http-api/api/v1/channels/delete"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/channels/:channelId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		channelsDelete.Handler,
		"successfully deletes channel",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testChannelId, _ := id.Generator.NewChannelId()

			expectedDatabaseOutputChannel := db.ChannelModel{
				InnerChannel: db.InnerChannel{
					Config: types.JSON{'{', '}'}, // empty config - required by mock Prisma client
				},
			}

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindUnique(
					db.Channel.ID.Equals(*testChannelId),
				),
			).Returns(expectedDatabaseOutputChannel)

			mocks.DB.Server.Channel.Expect(
				mocks.DB.Client.Channel.FindUnique(
					db.Channel.ID.Equals(*testChannelId),
				).Delete(),
			).Returns(expectedDatabaseOutputChannel)

			// run integration test
			routeWithChannelId := strings.ReplaceAll(route, ":channelId", *testChannelId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithChannelId, *testSpaceId)
			httptestRequest := httptest.NewRequest("DELETE", routeWithQueryParams, nil)
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, err := mocks.API.App.Test(httptestRequest)

			// assert nil error
			assert.Nil(t, err)

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		},
	)
}
