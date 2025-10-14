package search_test

import (
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
	"github.com/iotea-com/iotea/services/http-api/api/v1/users/search"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/users/search"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		search.Handler,
		"successfully searches users",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// define prisma mock server response
			testUserEmail := "test@iotea.com"
			testOrgId, _ := id.Generator.NewOrganizationId()

			expectedDatabaseOutput := []db.UserModel{
				{
					InnerUser: db.InnerUser{
						Email: &testUserEmail,
					},
				},
			}
			mocks.DB.Server.User.Expect(
				mocks.DB.Client.User.FindMany(
					db.User.Email.Contains(testUserEmail),
					db.User.Organizations.None(
						db.OrganizationMember.OrganizationID.Equals(*testOrgId),
					),
				).OrderBy(
					db.User.Relevance_.Fields([]db.UserOrderByRelevanceFieldEnum{
						db.UserOrderByRelevanceFieldEnumEmail,
					}),
					db.User.Relevance_.Search("database"),
					db.User.Relevance_.Sort("desc"),
				),
			).ReturnsMany(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			routeWithEmailQuery := route + "?q=" + testUserEmail
			routeWithQueryParams := fmt.Sprintf("%s&orgId=%s", routeWithEmailQuery, *testOrgId)
			httptestRequest := httptest.NewRequest("GET", routeWithQueryParams, nil)

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

			// check that Data is a []db.UserModel struct (marshal into JSON and unmarshal as an []user)
			usersJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var users []db.UserModel
			err = json.Unmarshal(usersJson, &users)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into []db.UserModel: %v", err)
			}

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Len(t, users, 1, "user ID matches database output")
			assert.Equalf(t, testUserEmail, *users[0].InnerUser.Email, "user ID matches database output")
		},
	)

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		search.Handler,
		"returns 404 if no users found",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// define prisma mock server response
			testOrgId, _ := id.Generator.NewOrganizationId()
			testUserEmail := "test@iotea.com"

			mocks.DB.Server.User.Expect(
				mocks.DB.Client.User.FindMany(
					db.User.Email.Contains(testUserEmail),
					db.User.Organizations.None(
						db.OrganizationMember.OrganizationID.Equals(*testOrgId),
					),
				).OrderBy(
					db.User.Relevance_.Fields([]db.UserOrderByRelevanceFieldEnum{
						db.UserOrderByRelevanceFieldEnumEmail,
					}),
					db.User.Relevance_.Search("database"),
					db.User.Relevance_.Sort("desc"),
				),
			).Errors(db.ErrNotFound)

			// when
			// API route is called with valid input but no users are found
			routeWithEmailQuery := route + "?q=" + testUserEmail
			routeWithQueryParams := fmt.Sprintf("%s&orgId=%s", routeWithEmailQuery, *testOrgId)
			httptestRequest := httptest.NewRequest("GET", routeWithQueryParams, nil)

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// assert response status
			assert.Equalf(t, http.StatusBadRequest, httptestResponse.StatusCode, "matching response code")
		},
	)
}
