package organizationsMembersRemove_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	organizationsMembersRemove "github.com/iotea-com/iotea/services/http-api/api/v1/organizations/members/remove"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/organizations/members"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		organizationsMembersRemove.Handler,
		"successfully removes user",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectOrgApiKey(mocks.DB)

			// define prisma mock server response
			testOrganizationId, _ := id.Generator.NewOrganizationId()
			testUserId := "mock-user-id"
			expectedDatabaseOutput := db.OrganizationMemberModel{
				InnerOrganizationMember: db.InnerOrganizationMember{
					OrganizationID: *testOrganizationId,
					UserID:         testUserId,
				},
				RelationsOrganizationMember: db.RelationsOrganizationMember{
					User: &db.UserModel{},
				},
			}

			mocks.DB.Server.OrganizationMember.Expect(
				mocks.DB.Client.OrganizationMember.FindUnique(
					db.OrganizationMember.OrganizationIDUserID(
						db.OrganizationMember.OrganizationID.Equals(*testOrganizationId),
						db.OrganizationMember.UserID.Equals(testUserId),
					),
				).Delete(),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"userId": testUserId,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithQueryParams := fmt.Sprintf("%s?orgId=%s", route, *testOrganizationId)
			httptestRequest := httptest.NewRequest("DELETE", routeWithQueryParams, bytes.NewBuffer(marshalledHttptestRequestBody))
			httptestRequest.Header.Add("Authorization", "Bearer tea_fake")

			httptestResponse, err := mocks.API.App.Test(httptestRequest)

			// then
			// assert nil error
			assert.Nil(t, err)

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
		},
	)
}
