package channelId

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/iotea-com/iotea/libs/engine/channels"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things/resolve"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Method", request.Input.Method),
		attribute.String("request.Input.ChannelId", request.Input.ChannelId),
	)

	// get channel configuration
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Get channel")
	c, err := prisma.Client.Channel.FindUnique(
		db.Channel.ID.Equals(request.Input.ChannelId),
	).Exec(dbCtx)
	if err != nil {
		if err.Error() == "ErrNotFound" {
			err := fmt.Sprintf("no channel found with ID %s", request.Input.ChannelId)
			request.Span.SetAttributes(
				attribute.String("error.type", "database"),
				attribute.String("error.message", err),
			)
			dbSpan.End()

			return nil, fiber.NewError(fiber.StatusBadRequest, err)
		}

		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error getting channel from the database: %s", err)),
		)
		dbSpan.End()

		return nil, err
	}

	dbSpan.End()

	// get HTTP server dependency configuration
	var channelConfig channels.Channel
	if err := json.Unmarshal([]byte(c.Config), &channelConfig); err != nil {
		err := fmt.Errorf("unable to unmarshal channel: %v", err)
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		request.Span.End()

		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	sourceNode, err := func() (*channels.Node, error) {
		// find the channel's source node
		for _, n := range channelConfig.Nodes {
			if n.Metadata.Type == "source" {
				return &n, nil
			}
		}

		return nil, fmt.Errorf("source node not found")
	}()

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// find all dependenices
	thingsFindConfig := resolve.FindConfig{
		Metadata:     &sourceNode.Metadata,
		PrismaClient: prisma.Client,
	}

	dependencies, err := resolve.FindDependenciesInNode(&thingsFindConfig)
	if err != nil {
		errorMessage := fmt.Sprintf("could not resolve things in node %v, %v", sourceNode.Metadata.Name, err)
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", errorMessage),
		)
		request.Span.End()
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	// get the first HTTP server dependency in the list of dependencies
	httpServerDependency, err := func() (*things.HttpServer, error) {
		for _, d := range dependencies {
			if d.ThingCategory == things.HttpServerThingCategory.String() {
				// attempt to unmarshal as an HTTP server
				var server things.HttpServer
				err := json.Unmarshal(d.Attributes, &server)
				if err != nil {
					err = fmt.Errorf("thing with ID %s was not a valid HTTP server - %s", d.ID, err.Error())
					request.Span.SetAttributes(
						attribute.String("error.type", "database"),
						attribute.String("error.message", err.Error()),
					)
					return nil, err
				}

				return &server, nil
			}
		}

		return nil, fiber.NewError(fiber.StatusBadRequest, "HTTP server dependency was not found")
	}()

	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// trigger source node using HTTP server dependency host and port
	var proxiedRequestBodyReader io.Reader = nil
	var contentType = "application/json" // Default to JSON
	if request.Input.Method != fiber.MethodGet && request.Input.Method != fiber.MethodHead {
		if len(request.Input.BinaryData) > 0 {
			// Handle Binary Upload
			proxiedRequestBodyReader = bytes.NewReader(request.Input.BinaryData)
			contentType = request.Input.BinaryContentType
		} else {
			proxiedRequestBody, err := json.Marshal(request.Input.Body)
			if err != nil {
				request.Span.SetAttributes(
					attribute.Bool("requestBody.forward", false),
				)
			}

			if len(proxiedRequestBody) > 0 {
				proxiedRequestBodyReader = bytes.NewReader(proxiedRequestBody)
			}
		}
	}

	httpRequest, err := http.NewRequestWithContext(request.Context, request.Input.Method, fmt.Sprintf("http://%s:%d", httpServerDependency.Host, httpServerDependency.Port), proxiedRequestBodyReader)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "http_request"),
			attribute.String("error.message", fmt.Sprintf("error creating HTTP request: %s", err)),
		)
		response, _ := ioteahttp.NewErrorResponse([]any{err.Error()}).MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(response))
	}

	// Set the appropriate Content-Type header
	httpRequest.Header.Set("Content-Type", contentType)

	response, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		request.Span.RecordError(fmt.Errorf("error sending HTTP request: %s", err))
		response, _ := ioteahttp.NewErrorResponse([]any{err.Error()}).MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(response))
	}

	output := Output{
		response: response,
	}

	return &output, nil
}
