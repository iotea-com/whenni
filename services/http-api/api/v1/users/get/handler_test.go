package get_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	"github.com/ongruent/gruent/services/http-api/api/v1/users/get"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/users/:userId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		get.Handler,
		"successfully gets user",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// define prisma mock server response
			testUserId := "mock-user-id"
			expectedDatabaseOutput := db.UserModel{
				InnerUser: db.InnerUser{
					ID: testUserId,
				},
			}
			mocks.DB.Server.User.Expect(
				mocks.DB.Client.User.FindUnique(
					db.User.ID.Equals(testUserId),
				).Omit(
					db.User.Password.Field(),
				).With(
					db.User.Organizations.Fetch().With(
						db.OrganizationMember.Organization.Fetch().With(
							db.Organization.Spaces.Fetch(),
						),
					),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			routeWithParams := strings.ReplaceAll(route, ":userId", testUserId)
			httptestRequest := httptest.NewRequest("GET", routeWithParams, nil)

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

			// check that Data is a db.UserModel struct (marshal into JSON and unmarshal as a user)
			userJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var user db.UserModel
			err = json.Unmarshal(userJson, &user)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.UserModel: %v", err)
			}

			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			assert.Equalf(t, user.ID, testUserId, "user ID matches database output")
		},
	)

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		get.Handler,
		"returns 404 if no user found",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// define prisma mock server response
			testUserId := "mock-user-id"
			mocks.DB.Server.User.Expect(
				mocks.DB.Client.User.FindUnique(
					db.User.ID.Equals(testUserId),
				).Omit(
					db.User.Password.Field(),
				).With(
					db.User.Organizations.Fetch().With(
						db.OrganizationMember.Organization.Fetch().With(
							db.Organization.Spaces.Fetch(),
						),
					),
				),
			).Errors(db.ErrNotFound)

			// when
			// API route is called with valid input but no user is found
			routeWithParams := strings.ReplaceAll(route, ":userId", testUserId)
			httptestRequest := httptest.NewRequest("GET", routeWithParams, nil)

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// assert response status
			assert.Equalf(t, http.StatusBadRequest, httptestResponse.StatusCode, "matching response code")
		},
	)
}
