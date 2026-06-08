package apitest

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/config"
)

type handlerDbMocks struct {
	Client *db.PrismaClient
	Server *db.Mock
}

type handlerApiMocks struct {
	Route string
	App   *fiber.App
}

type handlerSecretsMocks struct {
	Client *secrets.MockClient
}

type HandlerMocks struct {
	DB      handlerDbMocks
	API     handlerApiMocks
	Secrets handlerSecretsMocks
}

type handlerTestFunc = func(t *testing.T, mocks *HandlerMocks)

func HandlerUnitTestWithSetup(t *testing.T, route string, handlerFunc func(*fiber.Ctx) error, testDescription string, testFunc handlerTestFunc) bool {
	return t.Run(testDescription, func(t *testing.T) {
		// setup
		// set ENVIRONMENT to test
		os.Setenv("ENVIRONMENT", "test")

		// create mock prisma client
		dbClient, dbServer, ensure := db.NewMock()

		// ensure the mock database server received requests from the handler
		defer func() {
			// only run ensure if there were expectations - there may be no
			// expectations if the test does not reach the execute() func
			if dbServer.Expectations != nil {
				if len(*dbServer.Expectations) > 0 {
					ensure(t)
				}
			}
		}()

		// create mock Fiber client with route
		app := fiber.New()
		app.All(route, handlerFunc)

		// create a mock Vault server and client
		secretsClient := secrets.NewMockClient()
		config.SecretsClient = secretsClient

		// use mock ID generator
		id.Generator = &id.MockStandardIdGenerator{}

		// run the test logic
		mocks := &HandlerMocks{
			DB: handlerDbMocks{
				Client: dbClient,
				Server: dbServer,
			},
			API: handlerApiMocks{
				Route: route,
				App:   app,
			},
			Secrets: handlerSecretsMocks{
				Client: secretsClient,
			},
		}

		testFunc(t, mocks)
	})
}
