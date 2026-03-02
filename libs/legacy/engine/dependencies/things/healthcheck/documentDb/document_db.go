package documentDbHealthcheck

import (
	"context"
	"fmt"
	"time"

	"github.com/iotea-com/iotea/libs/legacy/engine/dependencies/things"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoDbServer(attrs *things.MongoDbServer) error {
	// Connect to MongoDB
	databaseUrl := (func() string {
		if attrs.Protocol == "mongodb" {
			return fmt.Sprintf("%s://%s:%s@%s:%d/", attrs.Protocol, attrs.Username, attrs.Password, attrs.Host, attrs.Port)
		}

		return fmt.Sprintf("%s://%s:%s@%s/", attrs.Protocol, attrs.Username, attrs.Password, attrs.Host)
	})()
	clientOptions := options.Client().ApplyURI(databaseUrl)

	dbCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Ping the server
	client, err := mongo.Connect(dbCtx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to a MongoDB server at %s:%d: %s", attrs.Host, attrs.Port, err)
	}

	err = client.Ping(dbCtx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping the MongoDB server at %s:%d: %s", attrs.Host, attrs.Port, err)
	}

	return nil
}
