package template

import (
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Replaces all instances of supported template functions with dynamic values.
//
// Example: {uuid()} -> ”aa816e61-44d3-4e22-838f-a83461813bb9"
func (t *Template) Preprocess() {
	t.preprocessNewUuid()
	t.preprocessNewDate()
	t.preprocessNewTimestamp()
}

func (t *Template) preprocessNewUuid() {
	uuidRegex := regexp.MustCompile(`uuid\(\)`)

	// ReplaceAllStringFunc is used to ensure a different UUID is used for each occurrence
	t.Value = uuidRegex.ReplaceAllStringFunc(t.Value, func(s string) string {
		return uuid.NewString()
	})
}

func (t *Template) preprocessNewDate() {
	rfc3339String := time.Now().Format(time.RFC3339)

	dateRegex := regexp.MustCompile(`date\(\)`)
	t.Value = dateRegex.ReplaceAllString(t.Value, rfc3339String)
}

func (t *Template) preprocessNewTimestamp() {
	timestamp := time.Now().Unix()

	timestampRegex := regexp.MustCompile(`timestamp\(\)`)
	t.Value = timestampRegex.ReplaceAllString(t.Value, fmt.Sprintf("%d", timestamp))
}
