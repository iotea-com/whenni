package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/iotea-com/iotea/libs/engine/dependencies/certificates"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
)

type ThingsIds struct {
	// HTTP Servers
	IoteaBeeceptorServer string

	// MQTT Brokers
	DockerMqttBroker string
	K8sMqttBroker    string

	// MQTTS Brokers
	DockerMqttsBroker string
	K8sMqttsBroker    string

	// MQTT Publishers
	DockerMqttPublisher string
	K8sMqttPublisher    string

	// MQTTS Publishers
	DockerMqttsPublisher string
	K8sMqttsPublisher    string

	// MQTT Subscribers
	DockerMqttSubscriber string
	K8sMqttSubscriber    string

	// MQTTS Subscribers
	DockerMqttsSubscriber string
	K8sMqttsSubscriber    string

	// Kafka Clusters
	DockerKafkaCluster string
	K8sKafkaCluster    string

	// Kafka Producers
	DockerKafkaProducer string
	K8sKafkaProducer    string

	// Kafka Consumers
	DockerKafkaConsumer string
	K8sKafkaConsumer    string

	// NATS Servers
	DockerNatsServer string
	K8sNatsServer    string

	// NATS Clients
	DockerNatsClient string
	K8sNatsClient    string

	// InluxDB Databases
	DockerInfluxDBDatabase string
	K8sInfluxDBDatabase    string
}

// Global service variable
var SeededThingsIds ThingsIds

func (s *DatabaseSeeder) CreateThings(spaceId, userId string) error {
	// HTTP Server
	if err := s.createHttpServerThings(spaceId); err != nil {
		return err
	}

	// MQTT Brokers
	if err := s.createMqttBrokerThings(spaceId); err != nil {
		return err
	}

	// MQTT Subscribers
	if err := s.createMqttSubscriberThings(spaceId); err != nil {
		return err
	}

	// MQTT Publishers
	if err := s.createMqttPublisherThings(spaceId); err != nil {
		return err
	}

	// Kafka Clusters
	if err := s.createKafkaClusterThings(spaceId); err != nil {
		return err
	}

	// Kafka Producers
	if err := s.createKafkaProducerThings(spaceId); err != nil {
		return err
	}

	// Kafka Consumers
	if err := s.createKafkaConsumerThings(spaceId); err != nil {
		return err
	}

	// NATS Servers
	if err := s.createNatsServerThings(spaceId); err != nil {
		return err
	}

	// NATS Clients
	if err := s.createNatsClientThings(spaceId); err != nil {
		return err
	}

	// InfluxDB Databases
	if err := s.createInfluxDBDatabaseThings(spaceId); err != nil {
		return err
	}

	log.Printf("Successfully upserted things into the database")

	return nil
}

func (s *DatabaseSeeder) createHttpServerThings(spaceId string) (err error) {
	category := things.HttpServerThingCategory.String()

	server := things.HttpServer{
		Host:     "iotea.free.beeceptor.com",
		Port:     443,
		Protocol: "https",
		Paths: []string{
			"/pass",
			"/fail",
			"/test1",
			"/test2",
			"/test3",
		},
	}

	SeededThingsIds.IoteaBeeceptorServer, err = s.marshalAndUploadThing(spaceId, server, category, "IOTEA Beeceptor")
	if err != nil {
		return err
	}

	return nil
}

func (s *DatabaseSeeder) createMqttBrokerThings(spaceId string) (err error) {
	category := things.MqttBrokerThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sMqttBroker := things.MqttBroker{
			Host:     "",
			Port:     1883,
			Protocol: "mqtt",
		}

		SeededThingsIds.K8sMqttBroker, err = s.marshalAndUploadThing(spaceId, k8sMqttBroker, category, "MQTT Broker (kubernetes)")
		if err != nil {
			return err
		}

		k8sMqttsBroker := things.MqttBroker{
			Host:     "",
			Port:     8883,
			Protocol: "mqtts",
			CaCert: certificates.CaCert{
				CaCertPem: SeededCerts.K8sCaCert,
			},
		}

		SeededThingsIds.K8sMqttsBroker, err = s.marshalAndUploadThing(spaceId, k8sMqttsBroker, category, "MQTTS Broker (kubernetes)")
		if err != nil {
			return err
		}
	} else { // Docker
		dockerMqttBroker := things.MqttBroker{
			Host:     "mosquitto",
			Port:     1883,
			Protocol: "mqtt",
		}

		SeededThingsIds.DockerMqttBroker, err = s.marshalAndUploadThing(spaceId, dockerMqttBroker, category, "MQTT Broker (docker)")
		if err != nil {
			return err
		}

		dockerMqttsBroker := things.MqttBroker{
			Host:     "emqx",
			Port:     8883,
			Protocol: "mqtts",
			CaCert: certificates.CaCert{
				CaCertPem: SeededCerts.DockerCaCert,
			},
		}

		SeededThingsIds.DockerMqttsBroker, err = s.marshalAndUploadThing(spaceId, dockerMqttsBroker, category, "MQTTS Broker (docker)")
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *DatabaseSeeder) createMqttSubscriberThings(spaceId string) (err error) {
	category := things.MqttClientThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sMqttSubscriber := things.MqttClient{
			Broker:   SeededThingsIds.K8sMqttBroker,
			Username: "MqttSubK8s",
			Password: "pass",
			ClientId: "MqttSubK8s",
		}

		SeededThingsIds.K8sMqttSubscriber, err = s.marshalAndUploadThing(spaceId, k8sMqttSubscriber, category, "MQTT Subscriber (kubernetes)")
		if err != nil {
			return err
		}

		k8sMqttsSubscriber := things.MqttClient{
			Broker:        SeededThingsIds.K8sMqttsBroker,
			Username:      "MqttsSubK8s",
			Password:      "pass",
			ClientId:      "MqttsSubK8s",
			CertificateId: SeededCerts.K8sMqttsSubscriberKeyPairId,
		}

		SeededThingsIds.K8sMqttsSubscriber, err = s.marshalAndUploadThing(spaceId, k8sMqttsSubscriber, category, "MQTTS Subscriber (kubernetes)")
		if err != nil {
			return err
		}
	} else { // Docker
		dockerMqttSubscriber := things.MqttClient{
			Broker:   SeededThingsIds.DockerMqttBroker,
			Username: "MqttSubDocker",
			Password: "pass",
			ClientId: "MqttSubDocker",
		}

		SeededThingsIds.DockerMqttSubscriber, err = s.marshalAndUploadThing(spaceId, dockerMqttSubscriber, category, "MQTT Subscriber (docker)")
		if err != nil {
			return err
		}

		dockerMqttsSubscriber := things.MqttClient{
			Broker:        SeededThingsIds.DockerMqttsBroker,
			Username:      "MqttsSubDocker",
			Password:      "pass",
			ClientId:      "MqttsSubDocker",
			CertificateId: SeededCerts.DockerMqttsSubscriberKeyPairId,
		}

		SeededThingsIds.DockerMqttsSubscriber, err = s.marshalAndUploadThing(spaceId, dockerMqttsSubscriber, category, "MQTTS Subscriber (docker)")
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *DatabaseSeeder) createMqttPublisherThings(spaceId string) (err error) {
	category := things.MqttClientThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sMqttPublisher := things.MqttClient{
			Broker:   SeededThingsIds.K8sMqttBroker,
			Username: "MqttPubK8s",
			Password: "pass",
			ClientId: "MqttPubK8s",
		}

		SeededThingsIds.K8sMqttPublisher, err = s.marshalAndUploadThing(spaceId, k8sMqttPublisher, category, "MQTT Publisher (kubernetes)")
		if err != nil {
			return err
		}

		k8sMqttsPublisher := things.MqttClient{
			Broker:        SeededThingsIds.K8sMqttsBroker,
			Username:      "MqttsPubK8s",
			Password:      "pass",
			ClientId:      "MqttsPubK8s",
			CertificateId: SeededCerts.K8sMqttsPublisherKeyPairId,
		}

		SeededThingsIds.K8sMqttsPublisher, err = s.marshalAndUploadThing(spaceId, k8sMqttsPublisher, category, "MQTTS Publisher (kubernetes)")
		if err != nil {
			return err
		}
	} else { // Docker
		dockerMqttPublisher := things.MqttClient{
			Broker:   SeededThingsIds.DockerMqttBroker,
			Username: "MqttPubDocker",
			Password: "pass",
			ClientId: "MqttPubDocker",
		}

		SeededThingsIds.DockerMqttPublisher, err = s.marshalAndUploadThing(spaceId, dockerMqttPublisher, category, "MQTT Publisher (docker)")
		if err != nil {
			return err
		}

		dockerMqttsPublisher := things.MqttClient{
			Broker:        SeededThingsIds.DockerMqttBroker,
			Username:      "MqttsPubDocker",
			Password:      "pass",
			ClientId:      "MqttsPubDocker",
			CertificateId: SeededCerts.DockerMqttsPublisherKeyPairId,
		}

		SeededThingsIds.DockerMqttsPublisher, err = s.marshalAndUploadThing(spaceId, dockerMqttsPublisher, category, "MQTTS Publisher (docker)")
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *DatabaseSeeder) createKafkaClusterThings(spaceId string) (err error) {
	category := things.KafkaClusterThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sCluster := things.KafkaCluster{
			BootstrapServers: []string{},
			Topics:           []string{},
		}

		SeededThingsIds.K8sKafkaCluster, err = s.marshalAndUploadThing(spaceId, k8sCluster, category, "Kafka Cluster (kubernetes)")
		if err != nil {
			return err
		}
	} else { // Docker
		dockerCluster := things.KafkaCluster{
			BootstrapServers: []string{
				"kafka:9092",
			},
			Topics: []string{
				"test1",
				"test2",
				"test3",
			},
		}

		SeededThingsIds.DockerKafkaCluster, err = s.marshalAndUploadThing(spaceId, dockerCluster, category, "Kafka Cluster (docker)")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DatabaseSeeder) createKafkaProducerThings(spaceId string) (err error) {
	category := things.KafkaProducerThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sProducer := things.KafkaProducer{
			Cluster: SeededThingsIds.K8sKafkaCluster,
		}

		SeededThingsIds.K8sKafkaProducer, err = s.marshalAndUploadThing(spaceId, k8sProducer, category, "Kafka Producer (kubernetes)")
		if err != nil {
			return err
		}

	} else { // Docker

		dockerProducer := things.KafkaProducer{
			Cluster: SeededThingsIds.DockerKafkaCluster,
		}

		SeededThingsIds.DockerKafkaProducer, err = s.marshalAndUploadThing(spaceId, dockerProducer, category, "Kafka Producer (docker)")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DatabaseSeeder) createKafkaConsumerThings(spaceId string) (err error) {
	category := things.KafkaConsumerThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sConsumer := things.KafkaConsumer{
			Cluster:         SeededThingsIds.K8sKafkaCluster,
			GroupId:         "",
			AutoOffsetReset: "",
		}

		SeededThingsIds.K8sKafkaConsumer, err = s.marshalAndUploadThing(spaceId, k8sConsumer, category, "Kafka Consumer (kubernetes)")
		if err != nil {
			return err
		}
	} else { // Docker
		dockerConsumer := things.KafkaConsumer{
			Cluster:         SeededThingsIds.DockerKafkaCluster,
			GroupId:         "test-group-id",
			AutoOffsetReset: "earliest",
		}

		SeededThingsIds.DockerKafkaConsumer, err = s.marshalAndUploadThing(spaceId, dockerConsumer, category, "Kafka Consumer (docker)")
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *DatabaseSeeder) createNatsServerThings(spaceId string) (err error) {
	category := things.NatsServerThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sServer := things.NatsServer{
			Host:   "",
			Port:   4444,
			Topics: []string{},
		}

		SeededThingsIds.K8sNatsServer, err = s.marshalAndUploadThing(spaceId, k8sServer, category, "NATS Server (kubernetes)")
		if err != nil {
			return err
		}
	} else { // Docker
		dockerServer := things.NatsServer{
			Host: "nats-service",
			Port: 4222,
			Topics: []string{
				"test1",
				"test2",
				"test3",
			},
		}

		SeededThingsIds.DockerNatsServer, err = s.marshalAndUploadThing(spaceId, dockerServer, category, "NATS Server (docker)")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DatabaseSeeder) createNatsClientThings(spaceId string) (err error) {
	category := things.NatsClientThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sClient := things.NatsClient{
			Server: SeededThingsIds.K8sNatsServer,
		}

		SeededThingsIds.K8sNatsClient, err = s.marshalAndUploadThing(spaceId, k8sClient, category, "NATS Client (kubernetes)")
		if err != nil {
			return err
		}

	} else { // Docker
		dockerClient := things.NatsClient{
			Server: SeededThingsIds.DockerNatsServer,
		}

		SeededThingsIds.DockerNatsClient, err = s.marshalAndUploadThing(spaceId, dockerClient, category, "NATS Client (docker)")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DatabaseSeeder) createInfluxDBDatabaseThings(spaceId string) (err error) {
	category := things.InfluxDbDatabaseThingCategory.String()

	if s.K8sEnabled { // K8s
		k8sDatabase := things.InfluxDbDatabase{
			Host:     "",
			Port:     8086,
			Protocol: "http",
			Token:    "",
			OrgName:  "",
		}

		SeededThingsIds.K8sInfluxDBDatabase, err = s.marshalAndUploadThing(spaceId, k8sDatabase, category, "InfluxDB Database (kubernetes)")
		if err != nil {
			return err
		}
	} else { // Docker
		dockerDatabase := things.InfluxDbDatabase{
			Host:     "influxdb",
			Port:     8086,
			Protocol: "http",
			Token:    "my-secret-influxdb2-token",
			OrgName:  "iotea",
		}

		SeededThingsIds.DockerInfluxDBDatabase, err = s.marshalAndUploadThing(spaceId, dockerDatabase, category, "InfluxDB Database (docker)")
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DatabaseSeeder) marshalAndUploadThing(spaceId string, thing any, category, name string) (thingId string, err error) {
	// Marshal attributes
	jsonAttributes, err := json.Marshal(&thing)
	if err != nil {
		return thingId, fmt.Errorf("error marshaling thing %v, %v", name, err)
	}

	// Create a new thing ID
	genId, err := id.Generator.NewThingId()
	if err != nil {
		return thingId, fmt.Errorf("error generating a new Thing ID for %v: %s", name, err)
	}

	thingId = *genId

	// Create a new Thing
	thing, err = s.prismaClient.Thing.CreateOne(
		db.Thing.ID.Set(thingId),
		db.Thing.Name.Set(name),
		db.Thing.Attributes.Set(jsonAttributes),
		db.Thing.CreatedBy.Set(uuid.New().String()),
		db.Thing.UpdatedBy.Set(uuid.New().String()),
		db.Thing.ThingCategory.Set(category),
		db.Thing.Space.Link(
			db.Space.ID.Equals(spaceId),
		),
	).Exec(context.Background())
	if err != nil {
		return thingId, fmt.Errorf("error inserting thing %v of type %v into the database: %v", name, category, err)
	}

	return thingId, nil
}
