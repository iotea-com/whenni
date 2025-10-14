package policiesUpdate

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	mqttPolicies "github.com/iotea-com/iotea/libs/http/policies"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	type RequestBody struct {
		Policy mqttPolicies.Policy `json:"policy"`
		Revoke bool                `json:"revoke"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	spaceId := ctx.Query("spaceId")
	certificateId := ctx.Params("certificateId")
	requestSpan.SetAttributes(
		attribute.String("query.spaceId", spaceId),
		attribute.String("params.certificateId", certificateId),
	)

	bearerToken, err := ioteahttputil.ParseBearerToken(ctx.Get("Authorization"))
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_bearer_token"),
			attribute.String("error.message", err.Error()),
		)
		return nil, err
	}

	var requestBody RequestBody
	err = json.Unmarshal(ctx.Body(), &requestBody)
	if err != nil {
		requestSpan.SetAttributes(
			attribute.String("error.type", "parse_request_body"),
			attribute.String("error.message", err.Error()),
		)
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity)
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken:   *bearerToken,
			SpaceId:       spaceId,
			CertificateId: certificateId,
			Policy:        requestBody.Policy,
			Revoke:        requestBody.Revoke,
		},
	}

	return request, nil
}
