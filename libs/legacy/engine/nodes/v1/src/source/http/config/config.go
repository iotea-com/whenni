package httpSourceNodeConfig

import (
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
)

type HttpResponseInfo struct {
	ResponseCode int               `json:"responseCode" validate:"required,min=100,max=599"`
	ResponseBody []byte            `json:"responseBody"`
	Headers      map[string]string `json:"headers"`
}

type HttpSourceNodeConfig struct {
	ThingHttpServer things.HttpServer `json:"httpServer::thing" validate:"required"`
	Method          string            `json:"method" validate:"required,oneof=GET POST PUT DELETE HEAD PATCH"`
	ResponseTimeout int               `json:"responseTimeout" validate:"required,min=1,max=30"`
}
