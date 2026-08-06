package permissionsUpdate

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	gruenthttputil "github.com/ongruent/gruent/libs/http/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func parse(ctx *fiber.Ctx) (*gruenthttp.Request[Input], error) {
	type RequestBody struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}

	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	orgId := ctx.Query("orgId")
	spaceId := ctx.Query("spaceId")
	permissionSetId := ctx.Params("permissionSetId")
	requestSpan.SetAttributes(
		attribute.String("query.orgId", orgId),
		attribute.String("query.spaceId", spaceId),
		attribute.String("params.permissionSetId", permissionSetId),
	)

	bearerToken, err := gruenthttputil.ParseBearerToken(ctx.Get("Authorization"))
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

	request := &gruenthttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken:     *bearerToken,
			OrgId:           orgId,
			SpaceId:         spaceId,
			PermissionSetId: permissionSetId,
			Name:            requestBody.Name,
			Permissions:     requestBody.Permissions,
		},
	}

	return request, nil
}
