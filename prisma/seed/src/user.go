package seed

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserEmail    = "test@iotea.com"
	UserPassword = "iotea!"
)

func (s *DatabaseSeeder) GetOrCreateTestUser() (string, error) {
	userId, err := id.Generator.NewUserId()
	if err != nil {
		return "", fmt.Errorf("error generating user id: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %v", err)
	}

	user, err := s.prismaClient.User.UpsertOne(
		db.User.Email.Equals(UserEmail),
	).Create(
		db.User.ID.Set(*userId),
		db.User.Email.Set(UserEmail),
		db.User.Password.Set(string(hashedPassword)),
		db.User.EmailVerified.Set(time.Now()),
	).Update().Exec(context.Background())

	if err != nil {
		return "", fmt.Errorf("error upserting user: %v", err)
	}

	log.Printf("Successfully upserted a test user with id: %s", user.ID)

	return user.ID, nil
}
