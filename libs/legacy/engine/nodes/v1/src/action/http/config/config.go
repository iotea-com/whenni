package httpActionNodeConfig

import "github.com/iotea-com/iotea/libs/engine/dependencies/things"

type HttpActionNodeConfig struct {
	ThingHttpServer things.HttpServer `json:"httpServer::thing" validate:"required"`
	Path            string            `json:"path" validate:"required"`
	Method          string            `json:"method" validate:"required,oneof=GET POST PUT DELETE HEAD PATCH"`
	Headers         map[string]string `json:"headers"`
	Body            map[string]any    `json:"body"`
}
