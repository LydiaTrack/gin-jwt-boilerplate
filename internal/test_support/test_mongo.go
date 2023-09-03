package test_support

import (
	"context"
	"fmt"
	"gin-jwt-boilerplate/internal/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestWithMongo(t *testing.T) {
	context := context.Background()
	container, err := mongodb.StartContainer(context)
	if err != nil {
		panic(err)
	}

	// Start mongodb client
	endpoint, err := container.Endpoint(context, "mongodb")
	if err != nil {
		t.Error(fmt.Errorf("failed to get endpoint: %w", err))
	}

	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(endpoint))
	if err != nil {
		t.Fatal(fmt.Errorf("error creating mongo client: %w", err))
	}

	err = mongoClient.Connect(context)
	if err != nil {
		t.Fatal(fmt.Errorf("error connecting to mongo: %w", err))
	}

	err = mongoClient.Ping(context, nil)
	if err != nil {
		t.Fatal(fmt.Errorf("error pinging mongo: %w", err))
	}

	defer container.Terminate(context)
}
