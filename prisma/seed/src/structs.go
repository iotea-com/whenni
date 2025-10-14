package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/iotea-com/iotea/libs/id"
	"github.com/iotea-com/iotea/prisma/db"
)

func (s *DatabaseSeeder) CreateStructs(spaceId, userId string) error {
	structsFolder := filepath.Join(Conf.AssetsFolder, "structs")

	// Read the directory
	entries, err := os.ReadDir(structsFolder)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}
	log.Printf("Found %d entries in structs folder", len(entries))

	// Iterate over each file in the directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Read the schema from the file
		filePath := filepath.Join(structsFolder, entry.Name())
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %v: %v", filePath, err)
		}

		var structMap map[string]any
		if err := json.Unmarshal(fileContent, &structMap); err != nil {
			return fmt.Errorf("error unmarshaling JSON from file %v: %v", filePath, err)
		}

		jsonSchema, err := json.Marshal(structMap)
		if err != nil {
			return fmt.Errorf("error marshaling JSON schema for file %v: %v", filePath, err)
		}

		// Generate a human-readable name from the file name
		baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		humanReadableName := capitalizeWords(strings.ReplaceAll(baseName, "_", " "))

		// Create a new thing ID
		genId, _ := id.Generator.NewThingId()

		structId := *genId

		// Upsert into database
		_, err = s.prismaClient.Struct.UpsertOne(
			db.Struct.ID.Equals(structId),
		).Create(
			db.Struct.ID.Set(structId),
			db.Struct.Name.Set(humanReadableName),
			db.Struct.Schema.Set(jsonSchema),
			db.Struct.CreatedBy.Set(userId),
			db.Struct.UpdatedBy.Set(userId),
			db.Struct.Space.Link(
				db.Space.ID.Equals(spaceId),
			),
		).Update(
			db.Struct.Name.Set(humanReadableName),
			db.Struct.Schema.Set(jsonSchema),
			db.Struct.UpdatedBy.Set(userId),
			db.Struct.Space.Link(
				db.Space.ID.Equals(spaceId),
			),
		).Exec(context.Background())

		if err != nil {
			return fmt.Errorf("error upserting struct into the database: %s", err)
		}

		log.Printf("Successfully upserted struct into the database: %v", structId)
	}

	return nil
}
