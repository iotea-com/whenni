package permissionsRemove_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	permissionsRemove "github.com/ongruent/gruent/services/http-api/api/v1/permissions/remove"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/permissions"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		permissionsRemove.Handler,
		"successfully deletes org permission set",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectOrgApiKey(mocks.DB)

			// define prisma mock server response
			testPermissionSetId, _ := id.Generator.NewPermissionSetId()
			expectedDatabaseOutput := db.PermissionSetModel{}
			mocks.DB.Server.PermissionSet.Expect(
				mocks.DB.Client.PermissionSet.FindUnique(
					db.PermissionSet.ID.Equals(*testPermissionSetId),
				),
			).Returns(expectedDatabaseOutput)
			mocks.DB.Server.PermissionSet.Expect(
				mocks.DB.Client.PermissionSet.FindUnique(
					db.PermissionSet.ID.Equals(*testPermissionSetId),
				).Delete(),
			).Returns(expectedDatabaseOutput)

			// run integration test
			httptestRequestBody := map[string]any{
				"permissionSetId": *testPermissionSetId,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			testOrganizationId, _ := id.Generator.NewOrganizationId()
			routeWithQueryParams := fmt.Sprintf("%s?orgId=%s", route, *testOrganizationId)
			httptestRequest := httptest.NewRequest("DELETE", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, err := mocks.API.App.Test(httptestRequest)

			// assert nil error
			assert.Nil(t, err)

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		},
	)
}
