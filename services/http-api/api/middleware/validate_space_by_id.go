package middleware

import (
	"context"

	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
)

func ValidateSpaceById(ctx *fiber.Ctx) error {
	spaceId := ctx.Params("spaceId")

	dbCtx := context.Background()
	_, err := sqlc.Queries.GetSpace(dbCtx, spaceId)
	if err != nil {
		if err == pgx.ErrNoRows {
			errorResponse := ioteahttp.NewErrorResponse([]any{
				"space not found",
			})
			errorResponseJson, _ := errorResponse.MarshalJson()
			return fiber.NewError(fiber.StatusUnprocessableEntity, string(errorResponseJson))
		}

		errorResponse := ioteahttp.NewErrorResponse([]any{
			err.Error(),
		})
		errorResponseJson, _ := errorResponse.MarshalJson()
		return fiber.NewError(fiber.StatusInternalServerError, string(errorResponseJson))
	}

	return ctx.Next()
}
