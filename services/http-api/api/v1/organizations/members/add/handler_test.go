package organizationsMembersAdd_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	organizationsMembersAdd "github.com/iotea-com/iotea/services/http-api/api/v1/organizations/members/add"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/organizations/members"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		organizationsMembersAdd.Handler,
		"successfully adds user to space",
		func(t *testing.T, mocks *apitest.HandlerMocks,
		) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectOrgApiKey(mocks.DB)

			// define prisma mock server response
			testOrganizationId, _ := id.Generator.NewOrganizationId()
			testUserId := "mock-user-id"
			testPermissionSetId, _ := id.Generator.NewPermissionSetId()

			expectedDefaultPermissionSet := db.PermissionSetModel{
				InnerPermissionSet: db.InnerPermissionSet{
					ID: *testPermissionSetId,
				},
			}
			mocks.DB.Server.PermissionSet.Expect(
				mocks.DB.Client.PermissionSet.FindFirst(
					db.PermissionSet.And(
						db.PermissionSet.OrganizationID.Equals(*testOrganizationId),
						db.PermissionSet.SpaceID.IsNull(),
						db.PermissionSet.Name.Equals("Default"),
					),
				),
			).Returns(expectedDefaultPermissionSet)

			expectedDatabaseOutput := db.OrganizationMemberModel{}
			mocks.DB.Server.OrganizationMember.Expect(
				mocks.DB.Client.OrganizationMember.CreateOne(
					db.OrganizationMember.Organization.Link(
						db.Organization.ID.Equals(*testOrganizationId),
					),
					db.OrganizationMember.User.Link(
						db.User.ID.Equals(testUserId),
					),
					db.OrganizationMember.OrganizationPermissionSet.Link(
						db.PermissionSet.ID.Equals(expectedDefaultPermissionSet.ID),
					),
					db.OrganizationMember.Role.Set(db.OrganizationRoleMember),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"userId": testUserId,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithQueryParams := fmt.Sprintf("%s?orgId=%s", mocks.API.Route, *testOrganizationId)
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

			// check that Data is a add.ResponseBody struct (marshal into JSON and unmarshal as a user add response)
			responseBodyJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var responseBody organizationsMembersAdd.Output
			err = json.Unmarshal(responseBodyJson, &responseBody)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into add.Responsebody: %v", err)
			}

			// assert response status and body
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		})
}
