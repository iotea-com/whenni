package modelsGet_test

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
	modelsGet "github.com/iotea-com/iotea/services/http-api/api/v1/models/get"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/models/:modelId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		modelsGet.Handler,
		"successfully gets model",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testModelId, _ := id.Generator.NewModelId()
			expectedDatabaseOutput := db.ModelModel{
				InnerModel: db.InnerModel{
					ID:         *testModelId,
					Attributes: types.JSON{'{', '}'}, // empty schema - required by mock Prisma client
				},
			}
			mocks.DB.Server.Model.Expect(
				mocks.DB.Client.Model.FindUnique(
					db.Model.ID.Equals(*testModelId),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			routeWithModelId := strings.ReplaceAll(route, ":modelId", *testModelId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithModelId, *testSpaceId)
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

			// check that Data is a db.ModelModel (marshal into JSON and unmarshal as a db.ModelModel)
			responseBodyJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var m db.ModelModel
			err = json.Unmarshal(responseBodyJson, &m)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.ModelModel: %v", err)
			}

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Equalf(t, m.ID, *testModelId, "model ID matches database output")
		},
	)
}
