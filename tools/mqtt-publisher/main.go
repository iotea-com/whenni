package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	MessageCount = 1
)

type Response struct {
	Token string `json:"token"`
}

type Payload struct {
	Adc0 float32 `json:"adc0"`
	Adc1 float32 `json:"adc1"`
	Adc2 float32 `json:"adc2"`
	Adc3 float32 `json:"adc3"`
	Adc4 float32 `json:"adc4"`
	Adc5 float32 `json:"adc5"`
	Adc6 float32 `json:"adc6"`
}

func NewTLSConfig() *tls.Config {
	// import CA cert
	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile("assets/certs/ca.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// import client cert and key
	cert, err := tls.LoadX509KeyPair("assets/certs/publisher.crt", "assets/certs/publisher.key")
	if err != nil {
		panic(err)
	}

	// configure TLS
	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert, // does not matter since the server is configured to check this
		ClientCAs:          nil,
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{cert},
	}
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	// Get user flags
	tlsPtr := flag.Bool("tls", false, "Set TLS option (true/false)")
	topicPtr := flag.String("topic", "", "Topic")
	flag.Parse()
	tlsEnabled := *tlsPtr
	topic := *topicPtr
	if topic == "" {
		fmt.Printf("Provide a topic using --topic")
		os.Exit(1)
	}
	fmt.Println("TLS Enabled:", tlsEnabled)
	fmt.Println("Topic:", topic)

	// Create and set all the relevant client options
	opts := mqtt.NewClientOptions()

	// Generate a random 4-digit number for this client ID
	randomNumber := fmt.Sprintf("%04d", rand.Intn(10000))
	opts.SetClientID(fmt.Sprintf("publisher_%s", randomNumber))

	// Set the right connection string and configuration for TLS
	if tlsEnabled {
		tlsconfig := NewTLSConfig()
		opts.SetTLSConfig(tlsconfig)
		opts.AddBroker("tls://localhost:8883")
	} else {
		opts.AddBroker("localhost:1883")
	}

	opts.SetDefaultPublishHandler(f)

	// Create a new client and start a connection
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("error connecting: %s", token.Error()))
	}

	fmt.Println("Client connected, publishing messages...")

	i := 0
	for range time.Tick(time.Duration(1) * time.Second) {
		if i == MessageCount {
			break
		}
		payload := Payload{
			Adc0: 1.2,
			Adc1: 2.3,
		}
		text, err := json.Marshal(payload)
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal struct, error:%v", err))
		}
		c.Publish(topic, 0, false, text)
		i++
	}

	fmt.Println("Done. Exiting")
}
