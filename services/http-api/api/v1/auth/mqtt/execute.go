package mqtt

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	mqttPolicies "github.com/iotea-com/iotea/libs/http/policies"
	"github.com/iotea-com/iotea/libs/val"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")

	// Extract cert SN
	sn, err := ExtractSerialNumber(request.Input.CertString)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "extract_serial_number"),
			attribute.String("error.message", err.Error()),
		)
		// respond with 200 and "result": "deny" so the request is blocked by the broker
		request.FiberContext.Status(200).JSON(ResponseBody{
			Result: "deny",
			Errors: []string{
				err.Error(),
			},
		})
		return nil, nil
	}

	request.Span.SetAttributes(
		attribute.String("request.Input.Topic", request.Input.Topic),
		attribute.String("request.Input.Action", request.Input.Action),
		attribute.String("certificate.SerialNumber", *sn),
	)

	// Get certificate from the database
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get certificate")
	certificate, err := prisma.Client.Certificate.FindUnique(
		db.Certificate.ID.Equals(*sn),
	).Exec(dbCtx)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting certificate from the database: %s", err)),
		)
		dbSpan.End()
		return &Output{
			Result: "deny",
		}, err
	}

	dbSpan.End()

	// Check if the topic is prefixed with the space ID
	request.Span.SetAttributes(
		attribute.String("certificate.SpaceID", certificate.SpaceID),
	)
	prefix := fmt.Sprintf("%s/", certificate.SpaceID)
	if !strings.HasPrefix(request.Input.Topic, prefix) {
		request.Span.SetAttributes(
			attribute.String("error.type", "topic_prefix"),
			attribute.String("error.message", fmt.Sprintf("topic %s is not prefixed with the certificate's space ID %s", request.Input.Topic, certificate.SpaceID)),
		)
		return &Output{
			Result: "deny",
		}, nil
	}

	// Check if the certificate policy allows the action
	result := func() string {
		// Check if the certificate is revoked
		if certificate.Revoke {
			return "deny"
		}

		// Unmarshal policy
		policy := mqttPolicies.Policy{}
		err := json.Unmarshal(certificate.Policy, &policy)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "unmarshal_policy"),
				attribute.String("error.message", fmt.Sprintf("failed to unmarshal policy: %v", err)),
			)
			return "deny"
		}

		// Validate the policy
		v := validator.New()
		v.RegisterValidation("mqtt_topic", val.MqttTopic)

		if err := v.Struct(policy); err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "validate_policy"),
				attribute.String("error.message", err.Error()),
			)
			return "deny"
		}

		// Assert policy
		err = policy.Assert(request.Input.Action, request.Input.Topic)
		if err != nil {
			request.Span.SetAttributes(
				attribute.String("error.type", "assert_policy"),
				attribute.String("error.message", err.Error()),
			)
			return "deny"
		}

		return "allow"
	}()

	output := &Output{
		Result: result,
	}

	return output, nil
}

// ExtractSerialNumber takes a certificate string and returns the serial number
func ExtractSerialNumber(certStr string) (*string, error) {
	// Ensure the certificate string has newlines at appropriate places
	certStr = formatCertificate(certStr)

	// Check if the certificate string contains the PEM headers
	if !strings.Contains(certStr, "-----BEGIN CERTIFICATE-----") || !strings.Contains(certStr, "-----END CERTIFICATE-----") {
		certStr = addPEMHeaders(certStr)
	}

	// Decode the PEM encoded certificate
	block, _ := pem.Decode([]byte(certStr))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing the certificate (cert: %s)", certStr)
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v - (cert: %s)", err, certStr)
	}

	serialNumberBytes := cert.SerialNumber.FillBytes(make([]byte, (cert.SerialNumber.BitLen()+7)/8))
	sn := fmt.Sprintf("%x", serialNumberBytes)
	var segments []string
	for i := 0; i < len(sn); i += 2 {
		if i+2 <= len(sn) {
			segments = append(segments, sn[i:i+2])
		} else {
			segments = append(segments, sn[i:])
		}
	}

	// Join the segments with a colon
	snString := strings.Join(segments, ":")
	return &snString, nil
}

// addPEMHeaders adds PEM headers to a certificate string if they are not present
func addPEMHeaders(certStr string) string {
	return "-----BEGIN CERTIFICATE-----\n" + certStr + "\n-----END CERTIFICATE-----"
}

// formatCertificate ensures that the certificate string has newlines at appropriate places
func formatCertificate(certStr string) string {
	// Remove existing newlines
	certStr = strings.ReplaceAll(certStr, "\n", "")

	// Add newlines after every 64 characters
	var formattedCert strings.Builder
	for i := 0; i < len(certStr); i += 64 {
		if i+64 < len(certStr) {
			formattedCert.WriteString(certStr[i:i+64] + "\n")
		} else {
			formattedCert.WriteString(certStr[i:] + "\n")
		}
	}
	return formattedCert.String()
}
