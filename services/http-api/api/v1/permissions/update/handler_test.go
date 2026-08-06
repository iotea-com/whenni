package permissionsUpdate_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	permissionsUpdate "github.com/ongruent/gruent/services/http-api/api/v1/permissions/update"
	"github.com/ongruent/gruent/services/http-api/util"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/organizations/:organizationId/permissions/:permissionSetId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		permissionsUpdate.Handler,
		"successfully updates org permission set",
		func(t *testing.T, mocks *apitest.HandlerMocks,
		) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectOrgApiKey(mocks.DB)

			// define prisma mock server response
			testApiKeyId, _ := id.Generator.NewApiKeyId()
			testPermissionSetId, _ := id.Generator.NewPermissionSetId()
			testOrganizationId, _ := id.Generator.NewOrganizationId()
			testPermissions := []string{"permissions:list"}
			testName := "Updated Name"

			expectedDatabaseOutput := db.PermissionSetModel{
				InnerPermissionSet: db.InnerPermissionSet{
					ID:             *testPermissionSetId,
					Name:           testName,
					OrganizationID: *testOrganizationId,
					CreatedBy:      *testApiKeyId,
					UpdatedBy:      *testApiKeyId,
					Permissions:    testPermissions,
				},
			}

			mocks.DB.Server.PermissionSet.Expect(
				mocks.DB.Client.PermissionSet.FindUnique(
					db.PermissionSet.ID.Equals(*testPermissionSetId),
				),
			).Returns(expectedDatabaseOutput)

			mocks.DB.Server.PermissionSet.Expect(
				mocks.DB.Client.PermissionSet.FindUnique(
					db.PermissionSet.ID.Equals(*testPermissionSetId),
				).Update(
					db.PermissionSet.Name.Set(testName),
					db.PermissionSet.UpdatedBy.Set(*testApiKeyId),
					db.PermissionSet.UpdatedAt.Set(util.GetCurrentTime()),
					db.PermissionSet.Permissions.Set(testPermissions),
					db.PermissionSet.Space.Link(
						db.Space.ID.EqualsIfPresent(nil),
					),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"name":        testName,
				"permissions": testPermissions,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithPermissionSetId := strings.ReplaceAll(route, ":permissionSetId", *testPermissionSetId)
			routeWithQueryParams := fmt.Sprintf("%s?orgId=%s", routeWithPermissionSetId, *testOrganizationId)
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
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Equal(t, permissionSet.ID, *testPermissionSetId, "permissionSet has the same ID as the database output")
			assert.Equal(t, permissionSet.Name, testName, "permissionSet has the same name as the database output")
			assert.Equal(t, permissionSet.Permissions[0], testPermissions[0], "permissionSet has the same permissions as the database output")
		})
}
