package ioteahttp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	ioteapermissions "github.com/iotea-com/iotea/libs/http/permissions"
	ioteahttputil "github.com/iotea-com/iotea/libs/http/util"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/util"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Actor struct {
	Id                    string
	CombinedPermissionSet map[string]*db.PermissionSetModel
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
	PrismaClient *db.PrismaClient           `validate:"-"` // skip validation for PrismaClient
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
			params.PrismaClient,
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
		params.PrismaClient,
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

func (r *Request[I]) setActor(id string, combinedPermissionSet map[string]*db.PermissionSetModel) {
	cps := make(map[string]*db.PermissionSetModel)
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
	p *db.PrismaClient,
	span trace.Span,
	scopeId string,
	unprotected bool,
	enforceAdmin bool,
	namespace ioteapermissions.Namespace,
	action ioteapermissions.Action,
) (string, map[string]*db.PermissionSetModel, error) {
	span.AddEvent("recognized API key in Authorization header")

	// Check if admin permissions are enforced
	if enforceAdmin {
		return "", nil, fmt.Errorf("admin permissions are enforced - API keys are not allowed to do this action")
	}

	// Get API key from the database
	// TODO: cache the results of this query
	dbCtx := context.Background()
	apiKey, err := p.APIKey.FindUnique(
		db.APIKey.ID.Equals(token),
	).With(
		db.APIKey.Organization.Fetch().With(
			db.Organization.Spaces.Fetch(),
		),
		db.APIKey.OrganizationPermissionSet.Fetch(),
		db.APIKey.SpacePermissionSet.Fetch(),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
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
	combinedPermissionSet := make(map[string]*db.PermissionSetModel)
	organizationPermissionSet, ok := apiKey.OrganizationPermissionSet()
	if ok {
		span.SetAttributes(attribute.StringSlice("authorize.apiKeyOrganizationPermissionSet", organizationPermissionSet.Permissions))
		combinedPermissionSet[apiKey.OrganizationID] = organizationPermissionSet
	}

	spacePermissionSet, ok := apiKey.SpacePermissionSet()
	if ok {
		permissionSetSpaceId, ok := spacePermissionSet.SpaceID()
		if ok {
			span.SetAttributes(attribute.StringSlice("authorize.apiKeySpacePermissionSet", spacePermissionSet.Permissions))
			combinedPermissionSet[permissionSetSpaceId] = spacePermissionSet
		}
	}

	// Check permissions in the organization or space
	allow := checkPermissions(
		db.OrganizationRoleMember,
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
	p *db.PrismaClient,
	jwtSecret string,
	span trace.Span,
	scopeId string,
	unprotected bool,
	enforceAdmin bool,
	namespace ioteapermissions.Namespace,
	action ioteapermissions.Action,
) (string, map[string]*db.PermissionSetModel, error) {
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
		} else {
			return userId, nil, nil
		}
	}

	// Get user permissions from the database
	organizationMember, err := func() (*db.OrganizationMemberModel, error) {
		switch scopeIdType {
		case id.IdTypeOrganization:
			// TODO: cache the results of this query
			dbCtx := context.Background()
			organizationMember, err := p.OrganizationMember.FindFirst(
				db.OrganizationMember.UserID.Equals(userId),
				db.OrganizationMember.OrganizationID.Equals(scopeId),
			).With(
				db.OrganizationMember.Organization.Fetch().With(
					db.Organization.Spaces.Fetch(),
				),
				db.OrganizationMember.User.Fetch(),
				db.OrganizationMember.OrganizationPermissionSet.Fetch(),
				db.OrganizationMember.SpacePermissions.Fetch(),
			).Exec(dbCtx)

			if err != nil {
				return nil, err
			}

			// If the OrganizationPermissionSet is nil, get the default organization permission set
			// TODO: cache the results of this query
			if organizationMember.OrganizationPermissionSet() == nil {
				dbCtx := context.Background()
				defaultOrganizationPermissionSet, err := p.PermissionSet.FindFirst(
					db.PermissionSet.Name.Equals("Default"),
					db.PermissionSet.OrganizationID.Equals(organizationMember.OrganizationID),
				).Exec(dbCtx)
				if err != nil {
					return nil, err
				}

				organizationMember = &db.OrganizationMemberModel{
					InnerOrganizationMember: organizationMember.InnerOrganizationMember,
					RelationsOrganizationMember: db.RelationsOrganizationMember{
						OrganizationPermissionSet: defaultOrganizationPermissionSet,
						User:                      organizationMember.User(),
						SpacePermissions:          organizationMember.SpacePermissions(),
						Organization:              organizationMember.Organization(),
					},
				}
			}

			return organizationMember, nil
		case id.IdTypeSpace:
			// Get the space from the database
			// TODO: cache the results of this query
			dbCtx := context.Background()
			space, err := p.Space.FindUnique(
				db.Space.ID.Equals(scopeId),
			).With(
				db.Space.Organization.Fetch(),
			).Exec(dbCtx)
			if err != nil {
				return nil, err
			}

			if space == nil {
				return nil, fmt.Errorf("space %s not found", scopeId)
			}

			// Get the organization member from the database using the space's organization ID
			// TODO: cache the results of this query
			dbCtx = context.Background()
			organizationMember, err := p.OrganizationMember.FindFirst(
				db.OrganizationMember.UserID.Equals(userId),
				db.OrganizationMember.OrganizationID.Equals(space.OrganizationID),
			).With(
				db.OrganizationMember.Organization.Fetch().With(
					db.Organization.Spaces.Fetch(),
				),
				db.OrganizationMember.User.Fetch(),
				db.OrganizationMember.OrganizationPermissionSet.Fetch(),
				db.OrganizationMember.SpacePermissions.Fetch().With(
					db.MemberSpacePermission.PermissionSet.Fetch(),
				),
			).Exec(dbCtx)
			if err != nil {
				return nil, err
			}

			// If the SpacePermissions is nil, get the default space permission set
			// TODO: cache the results of this query
			if len(organizationMember.SpacePermissions()) <= 0 {
				dbCtx := context.Background()
				defaultSpacePermissionSet, err := p.PermissionSet.FindFirst(
					db.PermissionSet.Name.Equals("Default"),
					db.PermissionSet.SpaceID.Equals(space.ID),
				).Exec(dbCtx)
				if err != nil {
					return nil, err
				}

				organizationMember = &db.OrganizationMemberModel{
					InnerOrganizationMember: organizationMember.InnerOrganizationMember,
					RelationsOrganizationMember: db.RelationsOrganizationMember{
						User: organizationMember.User(),
						SpacePermissions: []db.MemberSpacePermissionModel{
							{
								InnerMemberSpacePermission: db.InnerMemberSpacePermission{
									ID:                   "Default",
									OrganizationMemberID: organizationMember.ID,
									PermissionSetID:      defaultSpacePermissionSet.ID,
									CreatedAt:            util.GetCurrentTime(),
								},
								RelationsMemberSpacePermission: db.RelationsMemberSpacePermission{
									PermissionSet: defaultSpacePermissionSet,
								},
							},
						},
						OrganizationPermissionSet: organizationMember.OrganizationPermissionSet(),
						Organization:              organizationMember.Organization(),
					},
				}
			}

			return organizationMember, nil
		}

		return nil, fmt.Errorf("invalid scope ID type: %s", scopeIdType)
	}()

	if err != nil {
		if err.Error() == "ErrNotFound" {
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
	combinedPermissionSet := make(map[string]*db.PermissionSetModel)
	organizationPermissionSet := organizationMember.OrganizationPermissionSet()
	if organizationPermissionSet != nil {
		combinedPermissionSet[organizationMember.OrganizationID] = organizationPermissionSet
	}

	spacePermissionSets := organizationMember.SpacePermissions()
	if len(spacePermissionSets) > 0 {
		for _, spacePermissionSet := range spacePermissionSets {
			permissionSet := spacePermissionSet.PermissionSet()
			spaceId, ok := permissionSet.SpaceID()
			if !ok {
				continue
			}

			combinedPermissionSet[spaceId] = permissionSet
		}
	}

	combinedPermissionSetJson, _ := json.Marshal(combinedPermissionSet)
	span.SetAttributes(attribute.String("authorize.combinedPermissionSet", string(combinedPermissionSetJson)))

	// Check if admin permissions are enforced
	if enforceAdmin {
		if organizationMember.Role != db.OrganizationRoleAdmin {
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
	role db.OrganizationRole,
	combinedPermissionSet map[string]*db.PermissionSetModel,
	scopeId string,
	namespace ioteapermissions.Namespace,
	action ioteapermissions.Action,
) bool {
	allow := false
	if role == db.OrganizationRoleAdmin {
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
