package environmentsDelete

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"go.opentelemetry.io/otel/trace"
)

type RequestBody struct {
	Region        string `json:"region"`
	EnvironmentID string `json:"environmentID"`
}

func parse(ctx *fiber.Ctx) (*ioteahttp.Request[Input], error) {
	requestSpan := trace.SpanFromContext(ctx.UserContext())
	requestSpan.AddEvent("parse")

	spaceId := ctx.Query("spaceId")

	bearerToken, err := ioteahttputil.ParseBearerToken(ctx.Get("Authorization"))
	if err != nil {
		return nil, err
	}

	var requestBody RequestBody
	err = json.Unmarshal(ctx.Body(), &requestBody)
	if err != nil {
		return nil, fiber.NewError(http.StatusUnprocessableEntity)
	}

	request := &ioteahttp.Request[Input]{
		Span:         requestSpan,
		FiberContext: ctx,
		Context:      ctx.UserContext(),
		Input: Input{
			BearerToken:   *bearerToken,
			SpaceId:       spaceId,
			Region:        requestBody.Region,
			EnvironmentID: requestBody.EnvironmentID,
		},
	}

	return request, nil
}
