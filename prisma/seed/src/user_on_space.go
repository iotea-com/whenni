package seed

import (
	"context"
	"fmt"
	"log"

	"github.com/iotea-com/iotea/prisma/db"
)

func (s *DatabaseSeeder) CreateUserOnSpace(userId, spaceId, permissionSetId string) error {
	log.Printf("Creating user on space with id: %s, space id: %s, permission set id: %s", userId, spaceId, permissionSetId)

	_, err := s.prismaClient.UserOnSpace.UpsertOne(
		db.UserOnSpace.UserIDSpaceID(
			db.UserOnSpace.UserID.Equals(userId),
			db.UserOnSpace.SpaceID.Equals(spaceId),
		),
	).Create(
		db.UserOnSpace.User.Link(
			db.User.ID.Equals(userId),
		),
		db.UserOnSpace.Space.Link(
			db.Space.ID.Equals(spaceId),
		),
		db.UserOnSpace.PermissionSet.Link(
			db.PermissionSet.ID.Equals(permissionSetId),
		),
		db.UserOnSpace.Role.Set(db.UserOnSpaceRoleOwner),
	).Update().Exec(context.Background())

	if err != nil {
		return fmt.Errorf("error upserting user on space: %v", err)
	}

	log.Println("Successfully upserted a test user to space relationship")

	return nil
}
