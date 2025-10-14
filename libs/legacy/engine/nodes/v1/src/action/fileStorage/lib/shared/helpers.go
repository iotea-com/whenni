package helpers

import (
	"bytes"
	"errors"
	"net/http"

	node "github.com/iotea-com/iotea/libs/engine/nodes/v1"
)

// Helper function for awsS3 and Minio nodes
func DetectDataContentType(ioData node.IoData) ([]byte, string, error) {
	var data []byte
	var contentType string

	switch v := ioData.Data.(type) {
	case []byte:
		data = v
		// The http.DetectContentType function uses predefined signatures to identify common file types based on their initial bytes.
		contentType = http.DetectContentType(data)

	case string:
		data = []byte(v)
		contentType = "application/json"

	default:
		return nil, "", errors.New("unexpected data type in IoData")
	}

	// Adjust the content type if DetectContentType gave a generic value
	switch contentType {
	case "application/octet-stream", "text/plain; charset=utf-8":
		// If content is JSON but didn't start with `{` or `[`, it might be detected as octet-stream or plain text
		if bytes.HasPrefix(data, []byte("{")) || bytes.HasPrefix(data, []byte("[")) {
			contentType = "application/json"
		}
	}

	return data, contentType, nil
}
