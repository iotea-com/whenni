package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"math/rand"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewTLSConfig() *tls.Config {
	// import CA cert
	certpool := x509.NewCertPool()
	pemCerts, err := os.ReadFile("assets/certs/ca.crt")
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	}

	// import client keypair
	cert, err := tls.LoadX509KeyPair("assets/certs/subscriber.crt", "assets/certs/subscriber.key")
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
	opts.SetClientID(fmt.Sprintf("subscriber_%s", randomNumber))

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
		panic(token.Error())
	}

	fmt.Println("Client connected and subscribed! Awaiting on messages...")
	c.Subscribe(topic, 0, nil)

	for {
		select {}
	}
}
