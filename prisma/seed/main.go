package main

import (
	"flag"
	"log"

	seed "github.com/iotea-com/iotea/prisma/seed/src"
)

func main() {
	// Define and parse the command-line argument
	k8sEnabled := flag.Bool("k8sEnabled", false, "Enable or disable K8s support (true or false)")
	flag.Parse()

	// Get a new package instance
	seeder := seed.New()

	// Set the seeder.K8sEnabled variable
	seeder.K8sEnabled = *k8sEnabled

	// Create a new user
	userId, err := seeder.GetOrCreateTestUser()
	if err != nil {
		log.Fatal(err)
	}

	// Create test space
	spaceId, err := seeder.CreateTestSpace(userId)
	if err != nil {
		log.Fatalf("Could not create space, %v", err)
	}

	// Create default owner permission set
	permissionSetId, err := seeder.CreateOwnerPermissionSet(userId, spaceId)
	if err != nil {
		log.Fatal(err)
	}

	// Create relation to test user profile and space
	err = seeder.CreateUserOnSpace(userId, spaceId, permissionSetId)
	if err != nil {
		log.Fatal(err)
	}

	// Create an API Key
	err = seeder.CreateApiKey(spaceId, permissionSetId)
	if err != nil {
		log.Fatal(err)
	}

	// Create structs
	err = seeder.CreateStructs(spaceId, userId)
	if err != nil {
		log.Fatal(err)
	}

	// Create certificates
	err = seeder.CreateCerts(spaceId, userId)
	if err != nil {
		log.Fatal(err)
	}

	// Create things
	err = seeder.CreateThings(spaceId, userId)
	if err != nil {
		log.Fatal(err)
	}

	// Create channels
	err = seeder.CreateChannels(spaceId, userId)
	if err != nil {
		log.Fatal(err)
	}
}
