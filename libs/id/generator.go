package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type IdGenerator interface {
	NewUserId() (*string, error)
	NewOrganizationId() (*string, error)
	NewSpaceId() (*string, error)
	NewModelId() (*string, error)
	NewThingId() (*string, error)
	NewNodeId() (*string, error)
	NewChannelId() (*string, error)
	NewApiKeyId() (*string, error)
	NewPermissionSetId() (*string, error)
	NewChannelExecutionId() (*string, error)
	NewTagId() (*string, error)
}

type StandardIdGenerator struct{}

type MockStandardIdGenerator struct{}

var Generator IdGenerator = &StandardIdGenerator{}

const StandardIdLength = 16

func (idg *StandardIdGenerator) generateId(prefix string) (*string, error) {
	nanoid, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", StandardIdLength-len(prefix))
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s%s", prefix, nanoid)

	return &id, nil
}

func (idg *StandardIdGenerator) NewAccountId() (*string, error) {
	return idg.generateId("a")
}

func (idg *StandardIdGenerator) NewUserId() (*string, error) {
	return idg.generateId("u")
}

func (idg *StandardIdGenerator) NewOrganizationId() (*string, error) {
	return idg.generateId("o")
}

func (idg *StandardIdGenerator) NewSpaceId() (*string, error) {
	return idg.generateId("s")
}

func (idg *StandardIdGenerator) NewModelId() (*string, error) {
	return idg.generateId("m")
}

func (idg *StandardIdGenerator) NewThingId() (*string, error) {
	return idg.generateId("t")
}

func (idg *StandardIdGenerator) NewNodeId() (*string, error) {
	return idg.generateId("n")
}

func (idg *StandardIdGenerator) NewChannelId() (*string, error) {
	return idg.generateId("c")
}

func (idg *StandardIdGenerator) NewApiKeyId() (*string, error) {
	// generate random bytes
	keyBytes := make([]byte, 21)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return nil, err
	}

	// encode
	encodedKeyBytes := hex.EncodeToString(keyBytes)
	apiKey := fmt.Sprintf("tea_%s", encodedKeyBytes)

	return &apiKey, nil
}

func (idg *StandardIdGenerator) NewPermissionSetId() (*string, error) {
	return idg.generateId("p")
}

func (idg *StandardIdGenerator) NewChannelExecutionId() (*string, error) {
	return idg.generateId("x")
}

func (idg *StandardIdGenerator) NewTagId() (*string, error) {
	return idg.generateId("g")
}

func (idg *MockStandardIdGenerator) generateId(prefix string) (*string, error) {
	nanoid := "gruentidstub"
	id := fmt.Sprintf("%s%s", prefix, nanoid)

	return &id, nil
}

func (idg *MockStandardIdGenerator) NewAccountId() (*string, error) {
	return idg.generateId("a")
}

func (idg *MockStandardIdGenerator) NewUserId() (*string, error) {
	return idg.generateId("u")
}

func (idg *MockStandardIdGenerator) NewOrganizationId() (*string, error) {
	return idg.generateId("o")
}

func (idg *MockStandardIdGenerator) NewSpaceId() (*string, error) {
	return idg.generateId("s")
}

func (idg *MockStandardIdGenerator) NewModelId() (*string, error) {
	return idg.generateId("m")
}

func (idg *MockStandardIdGenerator) NewThingId() (*string, error) {
	return idg.generateId("t")
}

func (idg *MockStandardIdGenerator) NewNodeId() (*string, error) {
	return idg.generateId("n")
}

func (idg *MockStandardIdGenerator) NewChannelId() (*string, error) {
	return idg.generateId("c")
}

func (idg *MockStandardIdGenerator) NewApiKeyId() (*string, error) {
	apiKey := "tea_fake"
	return &apiKey, nil
}

func (idg *MockStandardIdGenerator) NewPermissionSetId() (*string, error) {
	return idg.generateId("p")
}

func (idg *MockStandardIdGenerator) NewChannelExecutionId() (*string, error) {
	return idg.generateId("x")
}

func (idg *MockStandardIdGenerator) NewTagId() (*string, error) {
	return idg.generateId("g")
}
