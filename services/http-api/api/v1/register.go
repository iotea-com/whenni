package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/services/http-api/api/v1/apiKeys"
	"github.com/iotea-com/iotea/services/http-api/api/v1/auth"
	"github.com/iotea-com/iotea/services/http-api/api/v1/certificates"
	"github.com/iotea-com/iotea/services/http-api/api/v1/channels"
	"github.com/iotea-com/iotea/services/http-api/api/v1/models"
	"github.com/iotea-com/iotea/services/http-api/api/v1/organizations"
	"github.com/iotea-com/iotea/services/http-api/api/v1/permissions"
	"github.com/iotea-com/iotea/services/http-api/api/v1/policies"
	"github.com/iotea-com/iotea/services/http-api/api/v1/secrets"
	"github.com/iotea-com/iotea/services/http-api/api/v1/spaces"
	"github.com/iotea-com/iotea/services/http-api/api/v1/tags"
	"github.com/iotea-com/iotea/services/http-api/api/v1/things"
	trigger "github.com/iotea-com/iotea/services/http-api/api/v1/trigger/http"
	"github.com/iotea-com/iotea/services/http-api/api/v1/users"
)

func Register(app *fiber.App) {
	v1 := app.Group("/v1")

	// api keys
	apiKeys.Register(v1)

	// auth
	auth.Register(v1)

	// certificates
	certificates.Register(v1)

	// channels
	channels.Register(v1)

	// organizations
	organizations.Register(v1)

	// permissions
	permissions.Register(v1)

	// policies
	policies.Register(v1)

	// secrets
	secrets.Register(v1)

	// spaces
	spaces.Register(v1)

	// models
	models.Register(v1)

	// tags
	tags.Register(v1)

	// things
	things.Register(v1)

	// triggers
	trigger.Register(v1)

	// users
	users.Register(v1)
}
