package thingsCreate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	thingsCreate "github.com/ongruent/gruent/services/http-api/api/v1/things/create"
	"github.com/ongruent/gruent/services/http-api/util"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/things/:thingId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		thingsCreate.Handler,
		"successfully creates an HTTP server thing",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testThingId, _ := id.Generator.NewThingId()
			testAttributes := things.HttpServer{
				Host:     "gruent.com",
				Protocol: "https",
				Port:     443,
				Paths: []string{
					"/test",
				},
			}

			requestBody := map[string]any{
				"name":       "test",
				"category":   things.HttpServerThingCategory.String(),
				"attributes": testAttributes,
			}

			testAttributesJson, _ := json.Marshal(testAttributes)
			requestBodyJson, _ := json.Marshal(requestBody)

			expectedDatabaseThingOutput := db.ThingModel{
				InnerThing: db.InnerThing{
					Attributes: testAttributesJson,
				},
			}

			now := util.GetCurrentTime()
			mocks.DB.Server.Thing.Expect(
				mocks.DB.Client.Thing.CreateOne(
					db.Thing.ID.Set(*testThingId),
					db.Thing.Name.Set("test"),
					db.Thing.Attributes.Set(testAttributesJson),
					db.Thing.CreatedBy.Set("tea_fake"),
					db.Thing.UpdatedBy.Set("tea_fake"),
					db.Thing.ThingCategory.Set(things.HttpServerThingCategory.String()),
					db.Thing.Space.Link(
						db.Space.ID.Equals(*testSpaceId),
					),
					db.Thing.CreatedAt.Set(now),
					db.Thing.UpdatedAt.Set(now),
				),
			).Returns(expectedDatabaseThingOutput)

			// run integration test
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", route, *testSpaceId)
			httptestRequest := httptest.NewRequest("POST", routeWithQueryParams, bytes.NewReader(requestBodyJson))
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")
			httptestResponse, err := mocks.API.App.Test(httptestRequest)

			// assert nil error
			assert.Nil(t, err)

			// assert response status
			assert.Equalf(t, http.StatusCreated, httptestResponse.StatusCode, "matching response code")
		},
	)

	// TODO: create tests for MQTT broker and MQTT client types when external certificates are supported
}
