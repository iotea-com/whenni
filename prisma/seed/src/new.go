package seed

import (
	"log"
	"os"

	"github.com/iotea-com/iotea/libs/secrets"
	"github.com/iotea-com/iotea/prisma/db"
)

type DatabaseSeeder struct {
	K8sEnabled    bool
	prismaClient  *db.PrismaClient
	secretsClient secrets.SecretsClient
}

func New() *DatabaseSeeder {
	seeder := DatabaseSeeder{}

	// Connect to database
	os.Setenv("DATABASE_URL", Conf.DatabaseUrl)
	seeder.prismaClient = db.NewClient()
	if err := seeder.prismaClient.Prisma.Connect(); err != nil {
		log.Fatalf("Could not connect to database client, %v", err)
	}

	// Connect secrets engine
	secretsClient, err := secrets.NewClient(Conf.VaultAddress, Conf.VaultUsername, Conf.VaultPassword, "")
	if err != nil {
		log.Fatalf("failed to log in to Vault")
	}

	seeder.secretsClient = secretsClient

	return &seeder
}
