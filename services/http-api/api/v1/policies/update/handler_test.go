package policiesUpdate_test

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
	mqttPolicies "github.com/ongruent/gruent/libs/http/policies"
	"github.com/ongruent/gruent/libs/id"
	"github.com/ongruent/gruent/prisma/db"
	apitest "github.com/ongruent/gruent/services/http-api/api/test"
	policiesUpdate "github.com/ongruent/gruent/services/http-api/api/v1/policies/update"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/policies/:certificateId"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		policiesUpdate.Handler,
		"successfully updates policy",
		func(t *testing.T, mocks *apitest.HandlerMocks,
		) {
			// given
			// define prisma mock server response for bearer token
			apitest.SetupExpectSpaceApiKey(mocks.DB)

			// define prisma mock server response
			testCertificateId, _ := id.Generator.NewPermissionSetId()
			testSpaceId, _ := id.Generator.NewSpaceId()
			testPolicy := mqttPolicies.Policy{
				AllowedSubscriptionTopics: []string{},
				AllowedPublishTopics:      []string{},
			}

			testPolicyJson, _ := json.Marshal(testPolicy)

			expectedDatabaseOutput := db.CertificateModel{
				InnerCertificate: db.InnerCertificate{
					Policy: testPolicyJson,
					Revoke: false,
				},
			}

			mocks.DB.Server.Certificate.Expect(
				mocks.DB.Client.Certificate.FindUnique(
					db.Certificate.ID.Equals(*testCertificateId),
				).Update(
					db.Certificate.Policy.Set(testPolicyJson),
					db.Certificate.Revoke.Set(false),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"policy": testPolicy,
				"revoke": false,
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)
			routeWithCertId := strings.ReplaceAll(route, ":certificateId", *testCertificateId)
			routeWithQueryParams := fmt.Sprintf("%s?spaceId=%s", routeWithCertId, *testSpaceId)
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

			// check that Data is a db.CertificateModel struct (marshal into JSON and unmarshal as a certificate)
			certificateJson, err := json.Marshal(response.Data)
			if err != nil {
				t.Fatalf("Failed to marshal response data back to JSON: %v", err)
			}

			var certificate db.CertificateModel
			err = json.Unmarshal(certificateJson, &certificate)
			if err != nil {
				t.Fatalf("Failed to unmarshal data into db.CertificateModel: %v", err)
			}

			// assert response status and body
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")
			// assert.Equal(t, certificate.Policy, testPermissions[0], "permissionSet has the same permissions as the database output")
			assert.Equal(t, certificate.Revoke, false, "certificate has the same revoke status as the database output")
		})
}
