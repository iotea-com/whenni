package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/ongruent/gruent/libs/legacy/engine/environment"
	"github.com/ongruent/gruent/libs/telemetry"
	"github.com/ongruent/gruent/services/http-api/api/healthcheck"
	v1 "github.com/ongruent/gruent/services/http-api/api/v1"
	"github.com/ongruent/gruent/services/http-api/config"
	clickhouseService "github.com/ongruent/gruent/services/http-api/services/clickhouse"
	sqlcService "github.com/ongruent/gruent/services/http-api/services/sqlc"
)

const (
	// How much time to give pending requests to complete when app gets OS KILL or SIGTERM
	ShutdownTimeout = 2 * time.Second
)

type Api struct {
	HttpServer *fiber.App
}

func New() *Api {
	// Instantiate new web server
	prefork := (func() bool {
		switch config.Env {
		case environment.Development:
			return false
		case environment.Production:
			return true
		default:
			return false
		}
	})()

	httpServer := fiber.New(fiber.Config{
		AppName:           "http-api",
		Prefork:           prefork,
		ServerHeader:      "Fiber",
		EnablePrintRoutes: true,
	})

	// Add middleware
	httpServer.Use(recover.New())
	httpServer.Use(cors.New())
	httpServer.Use(telemetry.TracerMiddleware())

	// Register core API groups
	healthcheck.Register(httpServer)
	v1.Register(httpServer)

	// Register Swagger
	httpServer.Static("/docs", "./api/static/docs")

	httpServer.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/docs/swagger.yaml",
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// Return API
	api := Api{
		HttpServer: httpServer,
	}

	return &api
}

func (a *Api) Listen(port int) error {
	return a.HttpServer.Listen(fmt.Sprintf(":%v", port))
}

func (a *Api) Shutdown() (errors []error) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	if err := a.HttpServer.ShutdownWithContext(shutdownCtx); err != nil {
		errors = append(errors, fmt.Errorf("server forced to shutdown: %v", err))
	}

	// Cleanup all the objects initialized in config.go

	if sqlcService.Pool != nil {
		sqlcService.Pool.Close()
	}

	if err := clickhouseService.Conn.Close(); err != nil {
		errors = append(errors, fmt.Errorf("could not disconnect from the ClickHouse database: %v", err))
	}

	return errors
}
