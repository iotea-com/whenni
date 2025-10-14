package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/iotea-com/iotea/prisma/db"
)

func (s *DatabaseSeeder) CreateTestSpace(userId string) (string, error) {
	// Upsert space
	space, err := s.prismaClient.Space.UpsertOne(
		db.Space.ID.Equals(userId),
	).Create(
		db.Space.ID.Set(userId),
		db.Space.Name.Set("Test User's Space"),
		db.Space.CreatedBy.Set(userId),
		db.Space.UpdatedBy.Set(userId),
	).Update(
		db.Space.ID.Set(userId), // No updates needed in the current scenario
	).Exec(context.Background())

	if err != nil {
		return "", fmt.Errorf("error upserting test space: %v", err)
	}

	log.Printf("Successfully upserted a test space %v", space.ID)

	return space.ID, nil
}
