package mqtt_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mqttPolicies "github.com/iotea-com/iotea/libs/http/policies"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	apitest "github.com/iotea-com/iotea/services/http-api/api/test"
	mqttAuth "github.com/iotea-com/iotea/services/http-api/api/v1/auth/mqtt"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/steebchen/prisma-client-go/runtime/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	route := "/v1/auth/mqtt"

	apitest.HandlerUnitTestWithSetup(
		t,
		route,
		mqttAuth.Handler,
		"successfully authorizes",
		func(t *testing.T, mocks *apitest.HandlerMocks) {
			// define prisma mock server response
			testCertificateString := "MIID5zCCAc+gAwIBAgIQVBM4qWK53jN8ZCel1vnzpTANBgkqhkiG9w0BAQsFADBlMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ08xDzANBgNVBAcMBkRlbnZlcjEUMBIGA1UECgwLSU9URUEsIEluYy4xDjAMBgNVBAsMBUlPVEVBMRIwEAYDVQQDDAlpb3RlYS5jb20wHhcNMjQwNzEwMjAxMzM5WhcNMjUwNzEwMjAxMzM5WjBlMQswCQYDVQQGEwJVUzELMAkGA1UECBMCQ08xDzANBgNVBAcTBkRlbnZlcjEUMBIGA1UEChMLSW9UZWEsIEluYy4xDjAMBgNVBAsTBUlvVGVhMRIwEAYDVQQDEwlsb2NhbGhvc3QwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQpAf3NGaYOIbJw7Az/6IgHtwnxrbF2G7sG2D1bmAcZHh2QVuUsIC+xIV6OJdXErGa0J6+UX7yqu0vODFCvVgN+o14wXDAfBgNVHSMEGDAWgBS3Ui87HdrW483s40X2s9z4h1jv2TAaBgNVHREEEzARgglsb2NhbGhvc3SHBH8AAAEwHQYDVR0OBBYEFE5uwRulsHZUp+61cnMoagbMqE9NMA0GCSqGSIb3DQEBCwUAA4ICAQAQ27F8NK/XiOEnSamwDnET93ZFm07ARpl0B6hs2jHha2FNoCzJrG+gZOYI2nNg/LOzOeAyxeFtmUXDHNK++hnW9e/DrUaPn5H4gzyWaYdqYVwOhlmilwss5ypyOMKR2kLCsx4YKmgoM2RcOlgWVzLjfA3blA8l+pyNMA+ebAV+LWhKCJy29DBf8nX0oMeVdy86vf8qR2wadP1Ap84iK0z0UadcPKQtBkzBVy7Dnr94fRiEAdhg+lpeyBho1JBq/6NPnKqLjMJ/vB8ypZzMwHC0IKymD8yGsQXoegC+MZ/A6MK99QsTp+97tN2XrVFrrWtiaV2oCTgK4mSak2MiyWVeh4YfEOKyPlpDLQwk5NhCK2wyfqPuvGipM/LlaVD9e6b4K1pc4Q6MD2xiMU/U9Qk2ziunV70DAxjTlFm7u8b5D59JInB3LHYSebgvptOZvxwlKpNGWzjLXFF5EFf5UF+Br7bMuWZTcMpOGS6dmyLTvySevb5sUlQ8x7YQOIqVuaKOVdVLFEG16W8Q5jsSj3rHijPREAK/bul23G3AyybFcvG2wnEX+hqhifxnEJsCotI69OpCd0E74y2h2twX/AsrQUFAsV4jDLdlI6CxVRkRvXxky7tqWTMZoucFYrDQeUFfbF1zNet9L3B3xX5tVRG2wxX1Ek9RSzh/xVY4B9yLtw=="
			sn, _ := mqttAuth.ExtractSerialNumber(testCertificateString)

			testSpaceId, _ := id.Generator.NewSpaceId()

			testPolicy := mqttPolicies.Policy{
				AllowedPublishTopics:      []string{fmt.Sprintf("%s/test", *testSpaceId)},
				AllowedSubscriptionTopics: []string{},
			}

			marshalledTestPolicy, _ := json.Marshal(testPolicy)

			expectedDatabaseOutput := db.CertificateModel{
				InnerCertificate: db.InnerCertificate{
					ID:      *sn,
					Name:    "Test Device",
					SpaceID: *testSpaceId,
					Revoke:  false,
					Policy:  types.JSON(marshalledTestPolicy),
				},
			}

			mocks.DB.Server.Certificate.Expect(
				sqlc.Client.Certificate.FindUnique(
					db.Certificate.ID.Equals(*sn),
				),
			).Returns(expectedDatabaseOutput)

			// when
			// API route is called with valid input
			httptestRequestBody := map[string]any{
				"certString": testCertificateString,
				"action":     "publish",
				"topic":      fmt.Sprintf("%s/test", *testSpaceId),
			}
			marshalledHttptestRequestBody, _ := json.Marshal(&httptestRequestBody)

			httptestRequest := httptest.NewRequest("POST", route, bytes.NewBuffer(marshalledHttptestRequestBody))

			httptestResponse, _ := mocks.API.App.Test(httptestRequest)

			// then
			// assert response status
			assert.Equalf(t, http.StatusOK, httptestResponse.StatusCode, "matching response code")

			// parse response body
			responseBodyBytes, err := io.ReadAll(httptestResponse.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}

			type ResponseBody struct {
				Result string `json:"result"`
			}
			var response ResponseBody
			err = json.Unmarshal(responseBodyBytes, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response body: %v", err)
			}

			// assert response body
			assert.Equal(t, "allow", response.Result, "result == allow")
		},
	)
}
