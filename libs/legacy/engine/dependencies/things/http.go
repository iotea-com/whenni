package things

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/ongruent/gruent/libs/val"
)

const HttpServerThingCategory ThingCategory = "HTTP_SERVER"

type HttpServer struct {
	Host     string   `json:"host" validate:"required,hostname|ip"`
	Port     int      `json:"port" validate:"required"`
	Protocol string   `json:"protocol" validate:"required,oneof=http https"`
	Paths    []string `json:"paths" validate:"required,valid_paths"`
}

func NewHttpServerFromAttributes(attributes map[string]any) (*HttpServer, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var httpServer HttpServer
	err := json.Unmarshal(jsonAttributes, &httpServer)
	if err != nil {
		return nil, err
	}

	// Validate
	err = httpServer.Validate()
	if err != nil {
		return nil, err
	}

	return &httpServer, nil
}

func (h *HttpServer) Category() ThingCategory {
	return HttpServerThingCategory
}

func (h *HttpServer) Validate() error {
	v := validator.New()
	v.RegisterValidation("valid_paths", val.IsValidPaths)

	if err := v.Struct(h); err != nil {
		return err
	}

	return nil
}

func (h *HttpServer) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(h)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}
