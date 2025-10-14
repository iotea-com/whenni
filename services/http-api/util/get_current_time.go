package util

import (
	"os"
	"time"
)

func GetCurrentTime() time.Time {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "test" {
		time, _ := time.Parse("yyyy-mm-dd", "2023-05-29")
		return time
	}

	return time.Now()
}
