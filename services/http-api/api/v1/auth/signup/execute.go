package signup

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Email", request.Input.Email),
	)

	// Generate a new user ID
	userId, err := id.Generator.NewUserId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "id_generation"),
			attribute.String("error.message", fmt.Sprintf("error generating user ID: %s", err)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Error creating user. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Input.Password), bcrypt.DefaultCost)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "password_hashing"),
			attribute.String("error.message", fmt.Sprintf("error hashing password: %s", err)),
		)
		errResponse := ioteahttp.NewErrorResponse([]any{"Error creating user. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	// Create a new user
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Create user")
	_, err = prisma.Client.User.CreateOne(
		db.User.ID.Set(*userId),
		db.User.Email.Set(request.Input.Email),
		db.User.Password.Set(string(hashedPassword)),
	).Exec(dbCtx)

	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error creating user: %s", err)),
		)
		dbSpan.End()
		errResponse := ioteahttp.NewErrorResponse([]any{"Error creating user. Please try again."})
		marshalledErrResponse, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
	}

	dbSpan.End()

	// If an orgId is provided, create a new organization member with the default org permission set
	if request.Input.InviteToken != nil {
		// Parse the invite token
		type claims struct {
			jwt.RegisteredClaims
			OrgId string `json:"orgId"`
		}
		inviteToken, err := jwt.ParseWithClaims(*request.Input.InviteToken, &claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected method: %s", token.Header["alg"])
			}
			return []byte(config.VaultConf.JwtSecret), nil
		})
		if err != nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "jwt_parsing"),
				attribute.String("error.message", fmt.Sprintf("error parsing invite token: %s", err)),
			)
			dbSpan.End()
			errResponse := ioteahttp.NewErrorResponse([]any{"Error parsing invite token. A user has been created, but you must be manually added to the organization."})
			marshalledErrResponse, _ := errResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusConflict, string(marshalledErrResponse))
		}

		parsedClaims := inviteToken.Claims.(*claims)

		// Get the default org permission set
		orgPermissionSet, err := prisma.Client.PermissionSet.FindFirst(
			db.PermissionSet.OrganizationID.Equals(parsedClaims.OrgId),
			db.PermissionSet.Name.Equals("Default"),
		).Exec(request.Context)

		if err != nil || orgPermissionSet == nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("error getting default org permission set: %s", err)),
			)
			dbSpan.End()
			errResponse := ioteahttp.NewErrorResponse([]any{"Error adding new user to organization. Please add the user to the organization manually."})
			marshalledErrResponse, _ := errResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
		}

		// Create a new organization member with the default org permission set
		dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Create organization member")
		_, err = prisma.Client.OrganizationMember.CreateOne(
			db.OrganizationMember.Organization.Link(
				db.Organization.ID.Equals(parsedClaims.OrgId),
			),
			db.OrganizationMember.User.Link(
				db.User.ID.Equals(*userId),
			),
			db.OrganizationMember.OrganizationPermissionSet.Link(
				db.PermissionSet.ID.Equals(orgPermissionSet.ID),
			),
			db.OrganizationMember.Role.Set("MEMBER"),
		).Exec(dbCtx)

		if err != nil {
			dbSpan.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", fmt.Sprintf("error getting default org permission set: %s", err)),
			)
			dbSpan.End()
			errResponse := ioteahttp.NewErrorResponse([]any{"Error adding new user to organization. Please add the user to the organization manually."})
			marshalledErrResponse, _ := errResponse.MarshalJson()
			return nil, fiber.NewError(fiber.StatusInternalServerError, string(marshalledErrResponse))
		}
	}

	return nil, nil
}
