package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/iotea-com/iotea/prisma/db"
)

func (s *DatabaseSeeder) CreateApiKey(spaceId, permissionSetId string) error {

	keyValue := fmt.Sprintf("tea_%v", spaceId)

	// Upsert into database
	apiKey, err := s.prismaClient.APIKey.UpsertOne(
		db.APIKey.ID.Equals(keyValue),
	).Create(
		db.APIKey.ID.Set(keyValue),
		db.APIKey.CreatedBy.Set(spaceId),
		db.APIKey.Space.Link(db.Space.ID.Equals(spaceId)),
		db.APIKey.PermissionSet.Link(db.PermissionSet.ID.Equals(permissionSetId)),
	).Update(
		db.APIKey.ID.Set(keyValue),
		db.APIKey.CreatedBy.Set(spaceId),
		db.APIKey.Space.Link(db.Space.ID.Equals(spaceId)),
		db.APIKey.PermissionSet.Link(db.PermissionSet.ID.Equals(permissionSetId)),
	).Exec(context.Background())

	if err != nil {
		return fmt.Errorf("error upserting apiKey into the database: %s", err)
	}

	log.Printf("Successfully upserted API Key into the database: %v", apiKey.ID)

	return nil
}
