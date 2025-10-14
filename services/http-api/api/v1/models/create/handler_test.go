package modelsCreate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iotea-com/iotea/libs/engine/dependencies/models"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	modelsCreate "github.com/iotea-com/iotea/services/http-api/api/v1/models/create"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/models"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		modelsCreate.Handler,
		"successfully creates model",
		func(t *testing.T, mocks *apitest.HandlerMocks,
		) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testSpaceId, _ := id.Generator.NewSpaceId()
			testModelName := "Test Model"
			testModelId, _ := id.Generator.NewModelId()
			testModelAttributes := map[string]models.Attribute{
				"f0": {
					Id:   "f0",
					Key:  "testField",
					Type: "string",
				},
			}
			marshalledTestModelAttributes, _ := json.Marshal(testModelAttributes)

			expectedDatabaseModelOutput := db.ModelModel{
				InnerModel: db.InnerModel{
					ID:         *testModelId,
					Attributes: marshalledTestModelAttributes,
				},
			}

			mocks.DB.Server.Model.Expect(
				mocks.DB.Client.Model.CreateOne(
					db.Model.ID.Set(*testModelId),
					db.Model.Name.Set(testModelName),
					db.Model.Attributes.Set(marshalledTestModelAttributes),
					db.Model.CreatedBy.Set("tea_fake"),
					db.Model.UpdatedBy.Set("tea_fake"),
					db.Model.Space.Link(
						db.Space.ID.Equals(*testSpaceId),
					),
				),
			).Returns(expectedDatabaseModelOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"name":       testModelName,
				"attributes": testModelAttributes,
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

			// assert response status and body
			assert.Equalf(t, http.StatusCreated, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, *testModelId, m.ID, "model has the same ID as the database output")
		})

	// Future tests - TODO:
	// case for struct limit in a space
}
