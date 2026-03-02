package modelsUpdate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/models"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	modelsUpdate "github.com/iotea-com/iotea/services/http-api/api/v1/models/update"
	"github.com/iotea-com/iotea/services/http-api/util"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/models/:modelId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		modelsUpdate.Handler,
		"successfully updates model",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testModelId, _ := id.Generator.NewModelId()
			testModelName := "Test Model"
			testModelAttributes := map[string]models.Attribute{
				"f0": {
					Id:   "f0",
					Key:  "testField",
					Type: "string",
				},
			}
			marshalledTestModelAttributes, _ := json.Marshal(testModelAttributes)
			expectedDatabaseOutput := db.ModelModel{
				InnerModel: db.InnerModel{
					ID:         *testModelId,
					Name:       testModelName,
					Attributes: marshalledTestModelAttributes,
				},
			}

			mocks.DB.Server.Model.Expect(
				mocks.DB.Client.Model.FindUnique(
					db.Model.ID.Equals(*testModelId),
				).Update(
					db.Model.Name.Set(testModelName),
					db.Model.Attributes.Set(marshalledTestModelAttributes),
					db.Model.UpdatedBy.Set("tea_fake"),
					db.Model.UpdatedAt.Set(util.GetCurrentTime()),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called
			httptestRequestBody := map[string]any{
				"model": map[string]any{
					"id":         *testModelId,
					"name":       testModelName,
					"attributes": testModelAttributes,
				},
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithModelId := strings.ReplaceAll(route, ":modelId", *testModelId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithModelId, *testSpaceId)
			httptestRequest := httptest.NewRequest("PUT", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
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
				t.Fatalf("Failed to unmarshal response body: %v - %s", err, string(responseBodyBytes))
			}

			if len(response.Errors) > 0 {
				t.Fatalf("Error in response: %s", response.Errors[0])
			}

			// check that Data is a db.ModelModel (marshal into JSON and unmarshal as a db.ModelModel)
			responseBodyJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var s db.ModelModel
			err = json.Unmarshal(responseBodyJson, &s)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.ModelModel: %v", err)
			}

			// assert response status and body
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, s.ID, *testModelId, "model has the same model ID as the database output")
		},
	)
}
