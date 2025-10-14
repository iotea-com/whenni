package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-json"

	"github.com/agoda-com/opentelemetry-go/otelzerolog"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-playground/validator/v10"
	"github.com/iotea-com/iotea/libs/engine/dependencies/things"
	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

const (
	// This timeout should allow enough room for a TLS handshake
	ConnectionTimeout = 5 * time.Second
)

func (n *MqttSourceNode) Init(params node.InitParams) node.Error {
	startTime := time.Now()

	// Setup tracer, meter and logger
	if err := n.setupObservability(params.Obsv); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	// Set the notify channel
	n.notifyChannel = params.NotifyChannel

	// Setup the error channel
	n.errChan = make(chan error)

	// Setup configuration fields
	if err := n.setConfig(params.Config); err != nil {
		return node.Error{
			Type:   node.ValidationError,
			Reason: fmt.Sprintf("invalid configuration, %v", err),
		}
	}

	// Connect to the MQTT broker
	if err := n.connectToBroker(); err != nil {
		return node.Error{
			Type:   node.FatalError,
			Reason: err.Error(),
		}
	}

	n.logger.Debug().Msgf("%s Node Initialized OK in %vms", NodeName, time.Since(startTime).Milliseconds())

	return node.Error{
		Type:   node.NoError,
		Reason: "",
	}
}

// setConfig unmarshalls the passed config byte array into the node's configuration
// struct, and validates all the configuration fields
func (n *MqttSourceNode) setConfig(config []byte) error {
	if err := json.Unmarshal(config, &n.config); err != nil {
		return fmt.Errorf("invalid node configuration JSON passed: %v", err)
	}

	// Validate config
	v := validator.New()

	if err := v.Struct(n.config); err != nil {
		return err
	}

	// Validate broker
	broker, ok := n.config.ThingMqttClient.Broker.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid broker configuration: expected a map but got %T with value %v. Is the thing being resolved correctly?", n.config.ThingMqttClient.Broker, n.config.ThingMqttClient.Broker)
	}

	brokerJSON, err := json.Marshal(broker)
	if err != nil {
		return fmt.Errorf("failed to marshal broker configuration: %v", err)
	}

	var mqttBroker things.MqttBroker
	if err := json.Unmarshal(brokerJSON, &mqttBroker); err != nil {
		return fmt.Errorf("invalid broker configuration: %v", err)
	}

	if err := v.Struct(mqttBroker); err != nil {
		return err
	}

	if mqttBroker.Protocol == "mqtts" && len(mqttBroker.CaCert.CaCertPem) == 0 {
		return fmt.Errorf("CA certificate for broker was empty, but protocol is set to MQTTS")
	}

	// Set the broker
	n.broker = mqttBroker

	return nil
}

// setupObservability will create a tracer, a meter and a logger to be used
// internally by the node, based on the providers passed at initialization
func (n *MqttSourceNode) setupObservability(obsv node.Observability) error {
	n.tracer = obsv.TracerProvider.Tracer("MqttSourceNodeTracer")
	n.meter = obsv.MeterProvider.Meter("MqttSourceNodeMeter")

	hook := otelzerolog.NewHook(obsv.LoggerProvider)
	n.logger = obsv.Logger.With().Str("node", NodeName).Logger().Hook(hook)

	return nil
}

// connectToBroker will instantiate a new MQTT client with options based on
// the user's configuration, and then attempt to connect to the provided
// broker.
func (n *MqttSourceNode) connectToBroker() error {
	// Create MQTT client options
	clientOpts, err := n.setClientOptions()
	if err != nil {
		return fmt.Errorf("could not set client options: %v", err)
	}

	// Create a client & connect to broker
	n.client = mqtt.NewClient(clientOpts)
	connToken := n.client.Connect()
	connToken.WaitTimeout(ConnectionTimeout)

	if err := connToken.Error(); err != nil {
		return fmt.Errorf("could not connect to MQTT broker, %v", err)
	}

	// We've connected succesfully, subscribe to the right topic
	if n.client.IsConnected() {
		// Subscribe to topic and register callback
		subToken := n.client.Subscribe(n.config.Topic, byte(n.config.QoS), n.mqttSubscriptionCallback)
		subToken.WaitTimeout(SubscriptionTimeout)
		if err := subToken.Error(); err != nil {
			return fmt.Errorf("error subscribing to topic, %v", err)
		}
	} else {
		return fmt.Errorf("no connection error, but not connected to broker either")
	}

	n.logger.Debug().Msgf("Succesfully connected to broker at %v:%v and subscribed to topic %v", n.broker.Host, n.broker.Port, n.config.Topic)

	return nil
}

// setClientOptions will create the MQTT options to be passed to the instantiated
// client, based on the user's configuration
func (n *MqttSourceNode) setClientOptions() (*mqtt.ClientOptions, error) {
	// Set the scheme to the right string based on what the MQTT client library expects
	var scheme string
	if n.broker.Protocol == "mqtt" {
		scheme = "tcp"
	} else if n.broker.Protocol == "mqtts" {
		scheme = "ssl"
	}

	// Create the broker URL with the config info
	brokerUrl := fmt.Sprintf("%v://%v:%v", scheme, n.broker.Host, n.broker.Port)

	// Set the client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerUrl)
	if n.config.ThingMqttClient.ClientId != "" {
		opts.SetClientID(n.config.ThingMqttClient.ClientId)
	}
	if n.config.ThingMqttClient.Username != "" {
		opts.SetUsername(n.config.ThingMqttClient.Username)
	}
	if n.config.ThingMqttClient.Password != "" {
		opts.SetPassword(n.config.ThingMqttClient.Password)
	}

	// Set TLS options if applicable
	if n.broker.Protocol == "mqtts" {
		tlsConfig, err := n.getTLSConfig()
		if err != nil {
			return nil, err
		}

		opts.SetTLSConfig(tlsConfig)
	}

	return opts, nil
}

// getTLSConfig will create a new TLS configuration for secure communications
// based on the user's configuration
func (n *MqttSourceNode) getTLSConfig() (*tls.Config, error) {
	// Transform string literals back into consumable PEM format
	caCertPem := strings.Replace(n.broker.CaCert.CaCertPem, `\n`, "\n", -1)
	clientCertPem := strings.Replace(n.config.ThingMqttClient.Keypair.ClientCertPem, `\n`, "\n", -1)
	clientKeyPem := strings.Replace(n.config.ThingMqttClient.Keypair.ClientCertKey, `\n`, "\n", -1)

	certpool := x509.NewCertPool()
	if !certpool.AppendCertsFromPEM([]byte(caCertPem)) {
		return nil, fmt.Errorf("could not append CA certificate to pool")
	}

	cert, err := tls.X509KeyPair([]byte(clientCertPem), []byte(clientKeyPem))
	if err != nil {
		return nil, fmt.Errorf("could not set client certificate and key: %s", err)
	}

	// Configure TLS
	return &tls.Config{
		// RootCAs = certs used to verify server cert.
		RootCAs: certpool,
		// ClientAuth = whether to request cert from server.
		// Since the server is set up for SSL, this happens
		// anyways.
		ClientAuth: tls.NoClientCert,
		// ClientCAs = certs used to validate client cert.
		ClientCAs: nil,
		// InsecureSkipVerify = verify that cert contents
		// match server. IP matches what is in cert etc.
		InsecureSkipVerify: false,
		// Certificates = list of certs client sends to server.
		Certificates: []tls.Certificate{cert},
	}, nil
}
