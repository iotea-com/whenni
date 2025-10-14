package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
)

// !! WARNING: MUST match type in "github.com/iotea-com/iotea/services/http-api/api/v1/spaces/policies/update"
type Policy struct {
	AllowedSubscriptionTopics []string `json:"allowedSubscriptionTopics" validate:"required"`
	AllowedPublishTopics      []string `json:"allowedPublishTopics" validate:"required"`
}

type CertsIds struct {
	// CA cert is still a string
	DockerCaCert string
	K8sCaCert    string

	// Subscribers
	DockerMqttsSubscriberKeyPairId string
	K8sMqttsSubscriberKeyPairId    string

	// Publisher
	DockerMqttsPublisherKeyPairId string
	K8sMqttsPublisherKeyPairId    string
}

// Global service variable
var SeededCerts CertsIds

func (s *DatabaseSeeder) CreateCerts(spaceId, userId string) error {
	// Get the root certificate as a raw string
	caCert, err := s.secretsClient.RetrieveRootCaCert()
	if err != nil {
		return fmt.Errorf("could not retrieve root CA certificate, %v", err)
	}

	// Create a map of client variable names to their respective pointers
	var clients map[string]*string
	if s.K8sEnabled {
		clients = map[string]*string{ // K8s
			"MQTTS Publisher (kubernetes)":  &SeededCerts.K8sMqttsPublisherKeyPairId,
			"MQTTS Subscriber (kubernetes)": &SeededCerts.K8sMqttsSubscriberKeyPairId,
		}

		SeededCerts.K8sCaCert = *caCert
	} else {
		clients = map[string]*string{ // Docker
			"MQTTS Publisher (docker)":  &SeededCerts.DockerMqttsPublisherKeyPairId,
			"MQTTS Subscriber (docker)": &SeededCerts.DockerMqttsSubscriberKeyPairId,
		}

		SeededCerts.DockerCaCert = *caCert
	}

	// Create certificates for each client and assign the serial number directly to the field
	for clientName, clientPtr := range clients {
		sn, err := s.createMqttClientCert(spaceId, clientName)
		if err != nil {
			return err
		}
		*clientPtr = sn // Dereference the pointer and assign the value
	}

	return nil
}

func (s *DatabaseSeeder) createMqttClientCert(spaceId, clientName string) (string, error) {
	// Request a new certificate
	certRequest := secrets.CertRequest{
		CommonName: "localhost",
		AltNames:   "localhost,emqx,emqx-mqtts.services.svc.cluster.local,mqtts.iotea.com",
		IpSans:     "127.0.0.1",
		Ttl:        "8760h",
		KeyType:    "ec",
		KeyBits:    256,
	}

	serialNumber, _, _, err := s.secretsClient.CreateCert(certRequest)
	if err != nil {
		return "", fmt.Errorf("error creating client certificate for: %s", err)
	}

	log.Printf("Created cert for %v with id %v", clientName, serialNumber)

	// Create a cert in the database
	defaultPolicy := Policy{
		AllowedSubscriptionTopics: []string{},
		AllowedPublishTopics:      []string{},
	}
	defaultPolicyJson, _ := json.Marshal(defaultPolicy)

	_, err = s.prismaClient.Certificate.CreateOne(
		db.Certificate.ID.Set(serialNumber),
		db.Certificate.Name.Set(clientName),
		db.Certificate.Policy.Set(defaultPolicyJson),
		db.Certificate.Space.Link(
			db.Space.ID.Equals(spaceId),
		),
	).Exec(context.Background())
	if err != nil {
		return "", fmt.Errorf("error inserting certificate into the database: %s", err)
	}

	return serialNumber, nil
}
