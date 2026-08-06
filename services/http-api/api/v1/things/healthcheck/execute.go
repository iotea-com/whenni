package thingsHealthcheck

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	gruenthttp "github.com/ongruent/gruent/libs/http"
	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
	documentDbHealthcheck "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/healthcheck/documentDb"
	fileStorageHealthcheck "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/healthcheck/fileStorage"
	httpHealthcheck "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/healthcheck/http"
	messageQueueHealthcheck "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/healthcheck/messageQueue"
	mqttHealthcheck "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/healthcheck/mqtt"
	notificationHealthcheck "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/healthcheck/notification"
	timeSeriesDbHealthcheck "github.com/ongruent/gruent/libs/legacy/engine/dependencies/things/healthcheck/timeSeriesDb"
	"github.com/ongruent/gruent/services/http-api/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func execute(request *gruenthttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.ThingCategory", request.Input.ThingCategory),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	switch request.Input.ThingCategory {
	case string(things.AwsSESThingCategory):
		// Init the thing attributes
		attrs, err := things.NewAwsSESFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = notificationHealthcheck.AwsSES(attrs)
		return nil, handleError(request.Span, err)
	case string(things.AwsSNSThingCategory):
		// Init the thing attributes
		attrs, err := things.NewAwsSNSFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = notificationHealthcheck.AwsSNS(attrs)
		return nil, handleError(request.Span, err)
	case string(things.HttpServerThingCategory):
		// Init the thing attributes
		attrs, err := things.NewHttpServerFromAttributes(request.Input.Attributes)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = httpHealthcheck.HttpServer(attrs)
		return nil, handleError(request.Span, err)
	case string(things.InfluxDbDatabaseThingCategory):
		// Init the thing attributes
		attrs, err := things.NewInfluxDbDatabaseFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = timeSeriesDbHealthcheck.InfluxDbDatabase(attrs)
		return nil, handleError(request.Span, err)
	case string(things.ClickhouseDatabaseThingCategory):
		// Init the thing attributes
		attrs, err := things.NewClickhouseDatabaseFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = timeSeriesDbHealthcheck.ClickhouseDatabase(attrs)
		return nil, handleError(request.Span, err)

	case string(things.KafkaClusterThingCategory):
		// Init the thing attributes
		attrs, err := things.NewKafkaClusterFromAttributes(request.Input.Attributes)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = messageQueueHealthcheck.KafkaCluster(attrs)
		return nil, handleError(request.Span, err)
	case string(things.KafkaConsumerThingCategory):
		// Init the thing attributes
		attrs, err := things.NewKafkaConsumerFromAttributes(request.Input.Attributes)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = messageQueueHealthcheck.KafkaConsumer(attrs)
		return nil, handleError(request.Span, err)
	case string(things.KafkaProducerThingCategory):
		// Init the thing attributes
		attrs, err := things.NewKafkaProducerFromAttributes(request.Input.Attributes)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = messageQueueHealthcheck.KafkaProducer(attrs)
		return nil, handleError(request.Span, err)
	case string(things.MinioBucketThingCategory):
		// Init the thing attributes
		attrs, err := things.NewMinioBucketFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = fileStorageHealthcheck.MinioBucket(attrs)
		return nil, handleError(request.Span, err)
	case string(things.MongoDbServerThingCategory):
		// Init the thing attributes
		attrs, err := things.NewMongoDbServerFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = documentDbHealthcheck.MongoDbServer(attrs)
		return nil, handleError(request.Span, err)
	case string(things.MqttBrokerThingCategory):
		// Init the thing attributes
		attrs, err := things.NewMqttBrokerFromAttributes(request.Input.Attributes)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = mqttHealthcheck.MqttBroker(attrs)
		return nil, handleError(request.Span, err)
	case string(things.MqttClientThingCategory):
		// Init the thing attributes
		attrs, err := things.NewMqttClientFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = mqttHealthcheck.MqttClient(attrs)
		return nil, handleError(request.Span, err)
	case string(things.NatsClientThingCategory):
		// Init the thing attributes
		attrs, err := things.NewNatsClientFromAttributes(request.Input.Attributes)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = messageQueueHealthcheck.NatsClient(attrs)
		return nil, handleError(request.Span, err)
	case string(things.NatsServerThingCategory):
		// Init the thing attributes
		attrs, err := things.NewNatsServerFromAttributes(request.Input.Attributes)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = messageQueueHealthcheck.NatsServer(attrs)
		return nil, handleError(request.Span, err)
	case string(things.SendgridClientThingCategory):
		// Init the thing attributes
		attrs, err := things.NewSendgridClientFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = notificationHealthcheck.SendgridClient(attrs)
		return nil, handleError(request.Span, err)
	case string(things.S3BucketThingCategory):
		// Init the thing attributes
		attrs, err := things.NewS3BucketFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
		if err != nil {
			return nil, handleError(request.Span, err)
		}

		// Run the healthcheck
		err = fileStorageHealthcheck.S3Bucket(attrs)
		return nil, handleError(request.Span, err)
	default:
		errResponse := gruenthttp.NewErrorResponse([]any{fmt.Sprintf("unrecognized thing category: %s", request.Input.ThingCategory)})
		response, _ := errResponse.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(response))
	}
}

func handleError(span trace.Span, err error) error {
	if err == nil {
		return nil
	}

	span.SetAttributes(attribute.String("error", err.Error()))
	errResponse := gruenthttp.NewErrorResponse([]any{err.Error()})
	response, _ := errResponse.MarshalJson()
	return fiber.NewError(fiber.StatusBadRequest, string(response))
}
