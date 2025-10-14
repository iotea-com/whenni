package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/iotea-com/iotea/prisma/db"
)

func (s *DatabaseSeeder) CreateOwnerPermissionSet(userId, spaceId string) (string, error) {
	ctx := context.Background()

	// Upsert owner permission set
	permissionSet, err := s.prismaClient.PermissionSet.UpsertOne(
		db.PermissionSet.ID.Equals(userId),
	).Create(
		db.PermissionSet.ID.Set(userId),
		db.PermissionSet.Name.Set("Owner"),
		db.PermissionSet.CreatedBy.Set(userId),
		db.PermissionSet.UpdatedBy.Set(userId),
		db.PermissionSet.Space.Link(
			db.Space.ID.Equals(spaceId),
		),
		db.PermissionSet.Permissions.Set([]string{"*"}),
	).Update(
		db.PermissionSet.Name.Set("Owner"),
		db.PermissionSet.UpdatedBy.Set(userId),
		db.PermissionSet.Space.Link(
			db.Space.ID.Equals(spaceId),
		),
		db.PermissionSet.Permissions.Set([]string{"*"}),
	).Exec(ctx)

	if err != nil {
		return "", fmt.Errorf("error upserting owner permission set: %v", err)
	}

	log.Printf("Successfully upserted an owner permission set for %v", spaceId)

	return permissionSet.ID, nil
}
