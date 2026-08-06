package permissionsAdd_test

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
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	permissionsAdd "github.com/ongruent/gruent/services/http-api/api/v1/permissions/add"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/organizations/:organizationId/permissions"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		permissionsAdd.Handler,
		"successfully inserts org permission set",
		func(t *testing.T, mocks *apitest.HandlerMocks,
		) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectOrgApiKey(mocks.DB)

			// define prisma mock server response
			testApiKeyId, _ := id.Generator.NewApiKeyId()
			testOrgId, _ := id.Generator.NewOrganizationId()
			testPermissions := []string{"*"}
			testName := "Test Permission Set"
			testPermissionSetId, _ := id.Generator.NewPermissionSetId()

			expectedDatabaseOutput := db.PermissionSetModel{
				InnerPermissionSet: db.InnerPermissionSet{
					ID:          *testPermissionSetId,
					Name:        testName,
					CreatedBy:   *testApiKeyId,
					UpdatedBy:   *testApiKeyId,
					Permissions: testPermissions,
				},
			}

			mocks.DB.Server.PermissionSet.Expect(
				mocks.DB.Client.PermissionSet.CreateOne(
					db.PermissionSet.ID.Set(*testPermissionSetId),
					db.PermissionSet.Name.Set(testName),
					db.PermissionSet.CreatedBy.Set(*testApiKeyId),
					db.PermissionSet.UpdatedBy.Set(*testApiKeyId),
					db.PermissionSet.Organization.Link(
						db.Organization.ID.Equals(*testOrgId),
					),
					db.PermissionSet.Space.Link(
						db.Space.ID.EqualsIfPresent(nil),
					),
					db.PermissionSet.Permissions.Set(testPermissions),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"name":        testName,
				"permissions": testPermissions,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithQueryParams := fmt.Sprintf("%s?orgId=%s", route, *testOrgId)
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

			// check that Data is a db.PermissionSetModel struct (marshal into JSON and unmarshal as a permission set)
			permissionSetJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var permissionSet db.PermissionSetModel
			err = json.Unmarshal(permissionSetJson, &permissionSet)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.PermissionSetModel: %v", err)
			}

			// assert response status and body
			assert.Equalf(t, http.StatusCreated, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, permissionSet.ID, *testPermissionSetId, "permissionSet has the same ID as the database output")
		})
}
