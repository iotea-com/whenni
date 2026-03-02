package thingsCreate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	ioteahttp "github.com/iotea-com/iotea/libs/http"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/certificates"
	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/things"

	mqttPolicies "github.com/iotea-com/iotea/libs/http/policies"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
	"github.com/iotea-com/iotea/services/http-api/config"
	"github.com/iotea-com/iotea/services/http-api/services/prisma"
	"github.com/iotea-com/iotea/services/http-api/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func execute(request *ioteahttp.Request[Input]) (*Output, error) {
	request.Span.AddEvent("execute")
	request.Span.SetAttributes(
		attribute.String("request.Input.Name", request.Input.Name),
		attribute.String("request.Input.Category", string(request.Input.Category)),
		attribute.String("request.Input.SpaceId", request.Input.SpaceId),
	)

	output := Output{
		Thing: nil,
	}

	thingId, err := id.Generator.NewThingId()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "generate_id"),
			attribute.String("error.message", fmt.Sprintf("error generating a new Thing ID: %s", err)),
		)
		return nil, err
	}

	switch request.Input.Category {
	case things.HttpServerThingCategory:
		thing, err := handleCreateHttpServer(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.MqttBrokerThingCategory:
		thing, err := handleCreateMqttBroker(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.MqttClientThingCategory:
		thing, err := handleCreateMqttClient(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.KafkaClusterThingCategory:
		thing, err := handleCreateKafkaCluster(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.KafkaProducerThingCategory:
		thing, err := handleCreateKafkaProducer(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.KafkaConsumerThingCategory:
		thing, err := handleCreateKafkaConsumer(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.NatsServerThingCategory:
		thing, err := handleCreateNatsServer(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.NatsClientThingCategory:
		thing, err := handleCreateNatsClient(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.S3BucketThingCategory:
		thing, err := handleCreateS3Bucket(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.MinioBucketThingCategory:
		thing, err := handleCreateMinioBucket(request, *thingId)
		if err != nil {
			return nil, err
		}
		output.Thing = thing
	case things.InfluxDbDatabaseThingCategory:
		thing, err := handleCreateInfluxDbDatabase(request, *thingId)
		if err != nil {
			return nil, err
		}

		output.Thing = thing
	case things.ClickhouseDatabaseThingCategory:
		thing, err := handleCreateClickhouseDatabase(request, *thingId)
		if err != nil {
			return nil, err
		}
		output.Thing = thing
	case things.AwsSESThingCategory:
		thing, err := handleCreateAwsSES(request, *thingId)
		if err != nil {
			return nil, err
		}
		output.Thing = thing
	case things.AwsSNSThingCategory:
		thing, err := handleCreateAwsSNS(request, *thingId)
		if err != nil {
			return nil, err
		}
		output.Thing = thing
	case things.SendgridClientThingCategory:
		thing, err := handleCreateSendgridClient(request, *thingId)
		if err != nil {
			return nil, err
		}
		output.Thing = thing
	case things.MongoDbServerThingCategory:
		thing, err := handleCreateMongoDbServer(request, *thingId)
		if err != nil {
			return nil, err
		}
		output.Thing = thing
	default:
		response := ioteahttp.NewErrorResponse([]any{
			"Invalid category.",
		})

		responseBody, _ := response.MarshalJson()
		return nil, fiber.NewError(fiber.StatusBadRequest, string(responseBody))
	}

	return &output, nil
}

func handleCreateMqttClient(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	mqttClient, err := things.NewMqttClientFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Request a new certificate
	certRequest := secrets.CertRequest{
		CommonName: "localhost",
		AltNames:   "localhost,emqx,emqx-mqtts.services.svc.cluster.local,mqtts.iotea.com",
		IpSans:     "127.0.0.1",
		Ttl:        "8760h",
		KeyType:    "ec",
		KeyBits:    256,
	}
	serialNumber, _, _, err := config.SecretsClient.CreateCert(certRequest)
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "create_certificate"),
			attribute.String("error.message", fmt.Sprintf("error creating certificate: %s", err)),
		)
		return nil, fmt.Errorf("error creating certificate: %s", err)
	}

	// Insert certificate into the database
	defaultPolicy := mqttPolicies.Policy{
		AllowedSubscriptionTopics: []string{},
		AllowedPublishTopics:      []string{},
	}
	defaultPolicyJson, _ := json.Marshal(defaultPolicy)

	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Create certificate")
	certificate, err := prisma.Client.Certificate.CreateOne(
		db.Certificate.ID.Set(serialNumber),
		db.Certificate.Name.Set(request.Input.Name),
		db.Certificate.Policy.Set(defaultPolicyJson),
		db.Certificate.Space.Link(
			db.Space.ID.Equals(request.Input.SpaceId),
		),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting certificate into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.SetAttributes(attribute.String("certificateId", certificate.ID))
	dbSpan.End()

	// Check if the MQTT Broker Thing exists
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Check MQTT Broker existence")
	_, err = prisma.Client.Thing.FindUnique(
		db.Thing.ID.Equals(mqttClient.Broker.(string)),
	).Exec(dbCtx)
	if err != nil {
		errMessage := fmt.Sprintf("error getting requested MQTT broker: %s", err)
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", errMessage),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	// Add certificate to attributes
	mqttClient.CertificateId = certificate.ID

	// Create a new Thing
	dbCtx, dbSpan = otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: mqttClient,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateMqttBroker(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	mqttBroker, err := things.NewMqttBrokerFromAttributes(request.Input.Attributes)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Add the CA certificate
	// TODO: support external MQTTS brokers
	certificate, err := config.SecretsClient.RetrieveRootCaCert()
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "secret_client"),
			attribute.String("error.message", err.Error()),
		)
		return nil, fiber.NewError(500, err.Error())
	}

	mqttBroker.CaCert = certificates.CaCert{
		CaCertPem: *certificate,
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: mqttBroker,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateHttpServer(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	httpServer, err := things.NewHttpServerFromAttributes(request.Input.Attributes)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: httpServer,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateKafkaCluster(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	kafkaCluster, err := things.NewKafkaClusterFromAttributes(request.Input.Attributes)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: kafkaCluster,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateKafkaProducer(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	kafkaProducer, err := things.NewKafkaProducerFromAttributes(request.Input.Attributes)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: kafkaProducer,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateKafkaConsumer(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	kafkaConsumer, err := things.NewKafkaConsumerFromAttributes(request.Input.Attributes)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: kafkaConsumer,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateNatsServer(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	natsServer, err := things.NewNatsServerFromAttributes(request.Input.Attributes)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: natsServer,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateNatsClient(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	natsClient, err := things.NewNatsClientFromAttributes(request.Input.Attributes)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: natsClient,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateS3Bucket(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	s3Bucket, err := things.NewS3BucketFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: s3Bucket,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateMinioBucket(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	minioBucket, err := things.NewMinioBucketFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: minioBucket,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateInfluxDbDatabase(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	influxDbDatabase, err := things.NewInfluxDbDatabaseFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: influxDbDatabase,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateClickhouseDatabase(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	clickhouseDatabase, err := things.NewClickhouseDatabaseFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: clickhouseDatabase,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateAwsSES(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	awsSES, err := things.NewAwsSESFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: awsSES,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateAwsSNS(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	awsSNS, err := things.NewAwsSNSFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: awsSNS,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateSendgridClient(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	sendgridClient, err := things.NewSendgridClientFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: sendgridClient,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

func handleCreateMongoDbServer(request *ioteahttp.Request[Input], thingId string) (*db.ThingModel, error) {
	// Create attributes
	mongoDbServer, err := things.NewMongoDbServerFromAttributes(request.Input.Attributes, config.SecretsClient, request.Input.SpaceId)
	err = handleCreateAttributesError(request, err)
	if err != nil {
		return nil, err
	}

	// Create a new Thing
	dbCtx, dbSpan := otel.Tracer("prisma").Start(request.Context, "Insert thing")
	thing, err := insertThingIntoDatabase(dbCtx, insertThingIntoDatabaseParams{
		SpaceId:    request.Input.SpaceId,
		ThingId:    thingId,
		ThingName:  request.Input.Name,
		ActorId:    request.GetActorId(),
		Attributes: mongoDbServer,
	})
	if err != nil {
		request.Span.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", err.Error()),
		)
		dbSpan.End()
		return nil, err
	}

	dbSpan.End()

	return thing, nil
}

type insertThingIntoDatabaseParams struct {
	SpaceId    string
	ThingId    string
	ThingName  string
	ActorId    string
	Attributes things.Attributes
}

func insertThingIntoDatabase(ctx context.Context, params insertThingIntoDatabaseParams) (*db.ThingModel, error) {
	jsonAttributes, err := params.Attributes.MarshalJson()
	if err != nil {
		return nil, err
	}

	now := util.GetCurrentTime()

	dbCtx, dbSpan := otel.Tracer("prisma").Start(ctx, "Insert thing")
	thing, err := prisma.Client.Thing.CreateOne(
		db.Thing.ID.Set(params.ThingId),
		db.Thing.Name.Set(params.ThingName),
		db.Thing.Attributes.Set(jsonAttributes),
		db.Thing.CreatedBy.Set(params.ActorId),
		db.Thing.UpdatedBy.Set(params.ActorId),
		db.Thing.ThingCategory.Set(params.Attributes.Category().String()),
		db.Thing.Space.Link(
			db.Space.ID.Equals(params.SpaceId),
		),
		db.Thing.CreatedAt.Set(now),
		db.Thing.UpdatedAt.Set(now),
	).Exec(dbCtx)
	if err != nil {
		dbSpan.SetAttributes(
			attribute.String("error.type", "database"),
			attribute.String("error.message", fmt.Sprintf("error inserting thing into the database: %s", err)),
		)
		dbSpan.End()
		return nil, err
	}

	return thing, nil
}

func handleCreateAttributesError(request *ioteahttp.Request[Input], err error) error {
	if err != nil {
		if err, ok := err.(validator.ValidationErrors); ok {
			request.Span.SetAttributes(attribute.String("validationError", err.Error()))
			validationError := ioteahttp.NewIoteaValidationError(err)
			response := validationError.MarshalResponse()
			responseJson, err := response.MarshalJson()
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return fiber.NewError(fiber.StatusBadRequest, string(responseJson))
		}
		request.Span.SetAttributes(
			attribute.String("error.type", "validation"),
			attribute.String("error.message", err.Error()),
		)
		return err
	}

	return nil
}
