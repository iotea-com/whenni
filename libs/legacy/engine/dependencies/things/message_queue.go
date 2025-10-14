package things

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

const KafkaClusterThingCategory ThingCategory = "KAFKA_CLUSTER"
const KafkaProducerThingCategory ThingCategory = "KAFKA_PRODUCER"
const KafkaConsumerThingCategory ThingCategory = "KAFKA_CONSUMER"

type KafkaCluster struct {
	BootstrapServers []string `json:"bootstrapServers" validate:"required"` // TODO: Add validations
	Topics           []string `json:"topics" validate:"required"`           // TODO: Add validations
}

func NewKafkaClusterFromAttributes(attributes map[string]any) (*KafkaCluster, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var kafkaCluster KafkaCluster
	err := json.Unmarshal(jsonAttributes, &kafkaCluster)
	if err != nil {
		return nil, err
	}

	// Validate
	err = kafkaCluster.Validate()
	if err != nil {
		return nil, err
	}

	return &kafkaCluster, nil
}

func (m *KafkaCluster) Category() ThingCategory {
	return KafkaClusterThingCategory
}

func (m *KafkaCluster) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *KafkaCluster) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

type KafkaProducer struct {
	Cluster any `json:"cluster" validate:"required"` // Should reference a Thing UUID or be a Thing attributes struct itself
}

func NewKafkaProducerFromAttributes(attributes map[string]any) (*KafkaProducer, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var kafkaProducer KafkaProducer
	err := json.Unmarshal(jsonAttributes, &kafkaProducer)
	if err != nil {
		return nil, err
	}

	// Validate
	err = kafkaProducer.Validate()
	if err != nil {
		return nil, err
	}

	return &kafkaProducer, nil
}

func (m *KafkaProducer) Category() ThingCategory {
	return KafkaProducerThingCategory
}

func (m *KafkaProducer) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *KafkaProducer) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

type KafkaConsumer struct {
	Cluster         any    `json:"cluster" validate:"required"` // Should reference a Thing UUID or be a Thing attributes struct itself
	GroupId         string `json:"groupId" validate:"required"`
	AutoOffsetReset string `json:"autoOffsetReset" validate:"required,oneof=latest earliest none"`
}

func NewKafkaConsumerFromAttributes(attributes map[string]any) (*KafkaConsumer, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var kafkaConsumer KafkaConsumer
	err := json.Unmarshal(jsonAttributes, &kafkaConsumer)
	if err != nil {
		return nil, err
	}

	// Validate
	err = kafkaConsumer.Validate()
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer, nil
}

func (m *KafkaConsumer) Category() ThingCategory {
	return KafkaConsumerThingCategory
}

func (m *KafkaConsumer) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *KafkaConsumer) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

const NatsServerThingCategory ThingCategory = "NATS_SERVER"
const NatsClientThingCategory ThingCategory = "NATS_CLIENT"

type NatsServer struct {
	Host   string   `json:"host" validate:"required,hostname|ip"`
	Port   int      `json:"port" validate:"required"`
	Topics []string `json:"topics" validate:"required"` // TODO: Add validations
}

func NewNatsServerFromAttributes(attributes map[string]any) (*NatsServer, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var natsServer NatsServer
	err := json.Unmarshal(jsonAttributes, &natsServer)
	if err != nil {
		return nil, err
	}

	// Validate
	err = natsServer.Validate()
	if err != nil {
		return nil, err
	}

	return &natsServer, nil
}

func (m *NatsServer) Category() ThingCategory {
	return NatsServerThingCategory
}

func (m *NatsServer) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *NatsServer) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}

type NatsClient struct {
	Server any `json:"server" validate:"required"` // Should reference a Thing UUID or be a Thing attributes struct itself
}

func NewNatsClientFromAttributes(attributes map[string]any) (*NatsClient, error) {
	// Unmarshal
	jsonAttributes, _ := json.Marshal(attributes)
	var natsClient NatsClient
	err := json.Unmarshal(jsonAttributes, &natsClient)
	if err != nil {
		return nil, err
	}

	// Validate
	err = natsClient.Validate()
	if err != nil {
		return nil, err
	}

	return &natsClient, nil
}

func (m *NatsClient) Category() ThingCategory {
	return NatsClientThingCategory
}

func (m *NatsClient) Validate() error {
	v := validator.New()

	if err := v.Struct(m); err != nil {
		return err
	}

	return nil
}

func (m *NatsClient) MarshalJson() ([]byte, error) {
	jsonAttributes, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return jsonAttributes, nil
}
