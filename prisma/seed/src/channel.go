package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/iotea-com/iotea/prisma/db"
)

func (s *DatabaseSeeder) CreateChannels(spaceId, userId string) error {
	channelsFolder := filepath.Join(Conf.AssetsFolder, "channels")

	// Read the directory
	entries, err := os.ReadDir(channelsFolder)
	if err != nil {
		return fmt.Errorf("error reading directory: %v", err)
	}
	log.Printf("Found %d entries in channels folder", len(entries))

	// Iterate over each file in the directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(channelsFolder, entry.Name())
		fileContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading file %v: %v", filePath, err)
		}

		var channelConfig map[string]any
		if err := json.Unmarshal(fileContent, &channelConfig); err != nil {
			return fmt.Errorf("error unmarshaling JSON from file %v: %v", filePath, err)
		}

		channelId, ok := channelConfig["id"].(string)
		if !ok {
			return fmt.Errorf("missing or invalid 'id' in file %v", filePath)
		}

		jsonConfig, err := json.Marshal(channelConfig)
		if err != nil {
			return fmt.Errorf("error marshaling JSON config for file %v: %v", filePath, err)
		}

		// Generate a human-readable name from the file name
		baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		humanReadableName := capitalizeWords(strings.ReplaceAll(baseName, "_", " "))

		// Upsert into database
		_, err = s.prismaClient.Channel.UpsertOne(
			db.Channel.ID.Equals(channelId),
		).Create(
			db.Channel.ID.Set(channelId),
			db.Channel.Name.Set(humanReadableName),
			db.Channel.Config.Set(jsonConfig),
			db.Channel.CreatedBy.Set(userId),
			db.Channel.UpdatedBy.Set(userId),
			db.Channel.Space.Link(
				db.Space.ID.Equals(spaceId),
			),
		).Update(
			db.Channel.Name.Set(humanReadableName),
			db.Channel.Config.Set(jsonConfig),
			db.Channel.UpdatedBy.Set(userId),
			db.Channel.Space.Link(
				db.Space.ID.Equals(spaceId),
			),
		).Exec(context.Background())

		if err != nil {
			return fmt.Errorf("error upserting channel into the database: %s", err)
		}

		log.Printf("Successfully upserted channel into the database: %v", channelId)
	}

	return nil
}

// capitalizeWords capitalizes the first letter of each word in a string
func capitalizeWords(s string) string {
	var result strings.Builder
	words := strings.Fields(s)
	for i, word := range words {
		if i > 0 {
			result.WriteRune(' ')
		}
		for j, r := range word {
			if j == 0 {
				result.WriteRune(unicode.ToUpper(r))
			} else {
				result.WriteRune(r)
			}
		}
	}
	return result.String()
}
