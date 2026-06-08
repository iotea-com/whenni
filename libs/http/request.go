package ioteahttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/libs/id"
	apisqlc "github.com/iotea-com/iotea/services/http-api/services/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PermissionSet struct {
	Permissions []string
}

type Actor struct {
	Id                    string
	CombinedPermissionSet map[string]*PermissionSet
}

type Request[I any] struct {
	Span         trace.Span
	FiberContext *fiber.Ctx
	Context      context.Context
	Input        I
	Actor        *Actor
}

var DefaultSpacePermissions = []ioteapermissions.Permission{
	"*",
}

type AuthorizeRequestParams struct {
	BearerToken  string                     `validate:"required"`
	JwtSecret    string                     `validate:"required"`
	ScopeId      string                     `validate:"required"` // ID of the user, organization, or space that the request is related to
	Unprotected  bool                       // If true, the request will not be checked for permissions
	EnforceAdmin bool                       // If true, the request will be checked for admin permissions
	Namespace    ioteapermissions.Namespace `validate:"required_if=Unprotected false"`
	Action       ioteapermissions.Action    `validate:"required_if=Unprotected false"`
}

func (r *Request[I]) Authorize(params AuthorizeRequestParams) *fiber.Error {
	v := validator.New()
	err := v.Struct(params)
	if err != nil {
		r.Span.SetAttributes(attribute.String("authorize.error", err.Error()))
		r.Span.AddEvent(fmt.Sprintf("could not validate the request's authorization parameters: %s", err))
		response := NewErrorResponse([]any{
			"could not validate the request's authorization parameters - please report this issue",
		})
		responseJson, _ := response.MarshalJson()
		return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
	}

	r.Span.SetAttributes(attribute.String("authorize.scopeId", params.ScopeId))

	// Handle API keys
	if strings.Contains(params.BearerToken, "tea_") {
		actorId, combinedPermissionSet, err := handleApiKey(
			params.BearerToken,
			r.Span,
			params.ScopeId,
			params.Unprotected,
			params.EnforceAdmin,
			params.Namespace,
			params.Action,
		)
		if err != nil {
			r.Span.AddEvent(err.Error())
			response := NewErrorResponse([]any{
				err.Error(),
			})
			responseJson, _ := response.MarshalJson()
			return fiber.NewError(fiber.StatusForbidden, string(responseJson))
		}

		r.setActor(actorId, combinedPermissionSet)
		return nil
	}

	// Handle JWTs
	actorId, combinedPermissionSet, err := handleJwt(
		params.BearerToken,
		params.JwtSecret,
		r.Span,
		params.ScopeId,
		params.Unprotected,
		params.EnforceAdmin,
		params.Namespace,
		params.Action,
	)
	if err != nil {
		r.Span.AddEvent(err.Error())
		response := NewErrorResponse([]any{
			err.Error(),
		})
		responseJson, _ := response.MarshalJson()
		return fiber.NewError(fiber.StatusForbidden, string(responseJson))
	}

	r.setActor(actorId, combinedPermissionSet)

	return nil
}

func (r *Request[I]) setActor(id string, combinedPermissionSet map[string]*PermissionSet) {
	cps := make(map[string]*PermissionSet)
	if combinedPermissionSet != nil {
		cps = combinedPermissionSet
	}

	r.Actor = &Actor{
		Id:                    id,
		CombinedPermissionSet: cps,
	}
}

func (r *Request[I]) GetActorId() string {
	if r.Actor == nil {
		return ""
	}

	return r.Actor.Id
}

func handleApiKey(
	token string,
	span trace.Span,
	scopeId string,
	unprotected bool,
	enforceAdmin bool,
	namespace ioteapermissions.Namespace,
	action ioteapermissions.Action,
) (string, map[string]*PermissionSet, error) {
	span.AddEvent("recognized API key in Authorization header")

	if apisqlc.Pool == nil {
		return "", nil, fmt.Errorf("database pool is not initialized")
	}

	// Check if admin permissions are enforced
	if enforceAdmin {
		return "", nil, fmt.Errorf("admin permissions are enforced - API keys are not allowed to do this action")
	}

	// Get API key from the database
	// TODO: cache the results of this query
	dbCtx := context.Background()
	var apiKey struct {
		ID                          string
		OrganizationID              string
		OrganizationPermissionSetID *string
		SpacePermissionSetID        *string
	}
	err := apisqlc.Pool.QueryRow(
		dbCtx,
		`SELECT id, organization_id, organization_permission_set_id, space_permission_set_id
		FROM app."apiKeys"
		WHERE id = $1`,
		token,
	).Scan(&apiKey.ID, &apiKey.OrganizationID, &apiKey.OrganizationPermissionSetID, &apiKey.SpacePermissionSetID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fmt.Errorf("no API key found in the database with ID: %s", token)
		}

		return "", nil, fmt.Errorf("error getting API key from the database: %s", err)
	}

	span.AddEvent("successfully got API key from the database")
	span.SetAttributes(attribute.String("authorize.apiKeyId", fmt.Sprintf("tea_*****%s", apiKey.ID[len(apiKey.ID)-5:])))
	span.SetAttributes(attribute.String("authorize.apiKeyOrganizationId", apiKey.OrganizationID))

	// Check if request is unprotected - we do this after getting the API key from the database
	// because we need to ensure the API key exists
	if unprotected {
		span.AddEvent("request is unprotected - skipping permissions check")
		return apiKey.ID, nil, nil
	}

	// Create combined permission map
	combinedPermissionSet := make(map[string]*PermissionSet)
	if apiKey.OrganizationPermissionSetID != nil {
		var organizationPermissions []string
		err = apisqlc.Pool.QueryRow(
			dbCtx,
			`SELECT permissions
			FROM app.permissions
			WHERE id = $1`,
			*apiKey.OrganizationPermissionSetID,
		).Scan(&organizationPermissions)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fmt.Errorf("error getting organization permission set from the database: %s", err)
		}
		organizationPermissionSet := &PermissionSet{Permissions: organizationPermissions}
		span.SetAttributes(attribute.StringSlice("authorize.apiKeyOrganizationPermissionSet", organizationPermissionSet.Permissions))
		combinedPermissionSet[apiKey.OrganizationID] = organizationPermissionSet
	}

	if apiKey.SpacePermissionSetID != nil {
		var permissionSetSpaceID *string
		var spacePermissions []string
		err = apisqlc.Pool.QueryRow(
			dbCtx,
			`SELECT space_id, permissions
			FROM app.permissions
			WHERE id = $1`,
			*apiKey.SpacePermissionSetID,
		).Scan(&permissionSetSpaceID, &spacePermissions)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fmt.Errorf("error getting space permission set from the database: %s", err)
		}
		if permissionSetSpaceID != nil {
			spacePermissionSet := &PermissionSet{Permissions: spacePermissions}
			span.SetAttributes(attribute.StringSlice("authorize.apiKeySpacePermissionSet", spacePermissionSet.Permissions))
			combinedPermissionSet[*permissionSetSpaceID] = spacePermissionSet
		}
	}

	// Check permissions in the organization or space
	allow := checkPermissions(
		"MEMBER",
		combinedPermissionSet,
		scopeId,
		namespace,
		action,
	)

	if !allow {
		errMessage := fmt.Sprintf("API key %s does not have permission to do '%s' action on '%s' namespace", apiKey.ID, action, namespace) // TODO: mask this value
		span.AddEvent(errMessage)
		return "", nil, fiber.NewError(fiber.StatusForbidden, errMessage)
	}

	span.AddEvent(fmt.Sprintf("API key has permission to do '%s' action on '%s' namespace", action, namespace))
	return apiKey.ID, combinedPermissionSet, nil
}

func handleJwt(
	token string,
	jwtSecret string,
	span trace.Span,
	scopeId string,
	unprotected bool,
	enforceAdmin bool,
	namespace ioteapermissions.Namespace,
	action ioteapermissions.Action,
) (string, map[string]*PermissionSet, error) {
	if apisqlc.Pool == nil {
		return "", nil, fmt.Errorf("database pool is not initialized")
	}

	// Get subject from JWT as user ID
	jwtToken, err := jwt.Parse(token, ioteahttputil.ParseJwt(jwtSecret))
	if err != nil {
		return "", nil, fmt.Errorf("could not parse JWT in authorization header: %s", err)
	}
	span.AddEvent("recognized JWT in Authorization header")
	userId, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return "", nil, fmt.Errorf("could not get user ID from JWT: %s", err)
	}

	span.SetAttributes(attribute.String("authorize.userId", userId))

	// Check if request is unprotected
	if unprotected {
		span.AddEvent("request is unprotected - skipping permissions check")
		return userId, nil, nil
	}

	// Parse scope ID
	scopeIdType, err := id.Parse(scopeId)
	if err != nil {
		span.AddEvent(err.Error())
		response := NewErrorResponse([]any{
			err.Error(),
		})
		responseJson, _ := response.MarshalJson()
		return "", nil, fiber.NewError(fiber.StatusUnauthorized, string(responseJson))
	}

	// If scope ID is a user ID, check if it matches the user ID in the JWT
	if scopeIdType == id.IdTypeUser {
		if scopeId != userId {
			errMessage := "user cannot make changes to another user"
			span.AddEvent(errMessage)
			response := NewErrorResponse([]any{
				errMessage,
			})
			responseJson, _ := response.MarshalJson()
			return "", nil, fiber.NewError(fiber.StatusForbidden, string(responseJson))
		}

		return userId, nil, nil
	}

	dbCtx := context.Background()
	organizationID := scopeId
	if scopeIdType == id.IdTypeSpace {
		err = apisqlc.Pool.QueryRow(
			dbCtx,
			`SELECT organization_id
			FROM app.spaces
			WHERE id = $1`,
			scopeId,
		).Scan(&organizationID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", nil, fmt.Errorf("space not found")
			}
			return "", nil, err
		}
	} else if scopeIdType != id.IdTypeOrganization {
		return "", nil, fmt.Errorf("invalid scope ID type: %s", scopeIdType)
	}

	// Get user permissions from the database
	var organizationMember struct {
		ID                          string
		OrganizationID              string
		Role                        string
		OrganizationPermissionSetID string
	}
	err = apisqlc.Pool.QueryRow(
		dbCtx,
		`SELECT id, organization_id, role, organization_permission_set_id
		FROM app.organization_members
		WHERE user_id = $1 AND organization_id = $2`,
		userId,
		organizationID,
	).Scan(
		&organizationMember.ID,
		&organizationMember.OrganizationID,
		&organizationMember.Role,
		&organizationMember.OrganizationPermissionSetID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			errMessage := "You are not a member of this organization or this space does not belong to the organization"
			span.AddEvent(errMessage)
			response := NewErrorResponse([]any{
				errMessage,
			})
			responseJson, _ := response.MarshalJson()
			return "", nil, fiber.NewError(fiber.StatusForbidden, string(responseJson))
		}

		if err.Error() == "space not found" {
			errMessage := fmt.Sprintf("Space with ID %s does not exist", scopeId)
			span.AddEvent(errMessage)
			response := NewErrorResponse([]any{
				errMessage,
			})
			responseJson, _ := response.MarshalJson()
			return "", nil, fiber.NewError(fiber.StatusForbidden, string(responseJson))
		}

		return "", nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Create combined permission map
	combinedPermissionSet := make(map[string]*PermissionSet)

	var organizationPermissions []string
	err = apisqlc.Pool.QueryRow(
		dbCtx,
		`SELECT permissions
		FROM app.permissions
		WHERE id = $1`,
		organizationMember.OrganizationPermissionSetID,
	).Scan(&organizationPermissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = apisqlc.Pool.QueryRow(
				dbCtx,
				`SELECT permissions
				FROM app.permissions
				WHERE organization_id = $1
				  AND space_id IS NULL
				  AND name = 'Default'
				LIMIT 1`,
				organizationMember.OrganizationID,
			).Scan(&organizationPermissions)
		}
		if err != nil {
			return "", nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	combinedPermissionSet[organizationMember.OrganizationID] = &PermissionSet{Permissions: organizationPermissions}

	spacePermissionRows, rowErr := apisqlc.Pool.Query(
		dbCtx,
		`SELECT p.space_id, p.permissions
		FROM app.member_space_permissions msp
		JOIN app.permissions p ON p.id = msp.permission_set_id
		WHERE msp.organization_member_id = $1`,
		organizationMember.ID,
	)
	if rowErr != nil {
		return "", nil, fiber.NewError(fiber.StatusInternalServerError, rowErr.Error())
	}
	for spacePermissionRows.Next() {
		var permissionSetSpaceID *string
		var permissions []string
		if scanErr := spacePermissionRows.Scan(&permissionSetSpaceID, &permissions); scanErr != nil {
			spacePermissionRows.Close()
			return "", nil, fiber.NewError(fiber.StatusInternalServerError, scanErr.Error())
		}
		if permissionSetSpaceID != nil {
			combinedPermissionSet[*permissionSetSpaceID] = &PermissionSet{Permissions: permissions}
		}
	}
	spacePermissionRows.Close()
	if spacePermissionRows.Err() != nil {
		return "", nil, fiber.NewError(fiber.StatusInternalServerError, spacePermissionRows.Err().Error())
	}

	if scopeIdType == id.IdTypeSpace && combinedPermissionSet[scopeId] == nil {
		var defaultSpacePermissions []string
		err = apisqlc.Pool.QueryRow(
			dbCtx,
			`SELECT permissions
			FROM app.permissions
			WHERE space_id = $1
			  AND name = 'Default'
			LIMIT 1`,
			scopeId,
		).Scan(&defaultSpacePermissions)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return "", nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if err == nil {
			combinedPermissionSet[scopeId] = &PermissionSet{Permissions: defaultSpacePermissions}
		}
	}

	combinedPermissionSetJson, _ := json.Marshal(combinedPermissionSet)
	span.SetAttributes(attribute.String("authorize.combinedPermissionSet", string(combinedPermissionSetJson)))

	// Check if admin permissions are enforced
	if enforceAdmin {
		if organizationMember.Role != "ADMIN" {
			errMessage := "admin permissions are enforced, and the user is not an admin"
			span.AddEvent(errMessage)
			response := NewErrorResponse([]any{
				errMessage,
			})
			responseJson, _ := response.MarshalJson()
			return "", nil, fiber.NewError(fiber.StatusForbidden, string(responseJson))
		}

		return userId, combinedPermissionSet, nil
	}

	// Check permissions in the organization or space
	allow := checkPermissions(
		organizationMember.Role,
		combinedPermissionSet,
		scopeId,
		namespace,
		action,
	)

	if !allow {
		errMessage := fmt.Sprintf("user %s does not have permission to do '%s' action on '%s' namespace", userId, action, namespace)
		span.AddEvent(errMessage)
		return "", nil, fiber.NewError(fiber.StatusForbidden, errMessage)
	}

	span.AddEvent(fmt.Sprintf("user %s has permission to do '%s' action on '%s' namespace", userId, action, namespace))
	return userId, combinedPermissionSet, nil
}

func checkPermissions(
	role string,
	combinedPermissionSet map[string]*PermissionSet,
	scopeId string,
	namespace ioteapermissions.Namespace,
	action ioteapermissions.Action,
) bool {
	allow := false
	if role == "ADMIN" {
		// bypass permissions check if role is admin
		allow = true
	} else {
		namespaceAll := fmt.Sprintf("%s:%s", namespace, ioteapermissions.ActionAll)
		namespaceAction := fmt.Sprintf("%s:%s", namespace, action)

		for sid, permissionSet := range combinedPermissionSet {
			if sid != scopeId {
				continue
			}

			for _, permission := range permissionSet.Permissions {
				switch permission {
				case "*", namespaceAll, namespaceAction:
					allow = true
				default:
					continue
				}
			}
		}
	}

	return allow
}
