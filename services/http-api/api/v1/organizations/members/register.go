package organizationsMembers

import (
	"github.com/gofiber/fiber/v2"

	organizationsMembersAdd "github.com/ongruent/gruent/services/http-api/api/v1/organizations/members/add"
	organizationsMembersChangeRole "github.com/ongruent/gruent/services/http-api/api/v1/organizations/members/changeRole"
	organizationsMembersInvite "github.com/ongruent/gruent/services/http-api/api/v1/organizations/members/invite"
	organizationsMembersList "github.com/ongruent/gruent/services/http-api/api/v1/organizations/members/list"
	organizationsMembersRemove "github.com/ongruent/gruent/services/http-api/api/v1/organizations/members/remove"
)

func Register(app fiber.Router) {
	api := app.Group("members")

	api.Get("/", organizationsMembersList.Handler)
	api.Post("/", organizationsMembersAdd.Handler)
	api.Post("/invite", organizationsMembersInvite.Handler)
	api.Delete("/", organizationsMembersRemove.Handler)
	api.Patch("/changeRole", organizationsMembersChangeRole.Handler)
}
