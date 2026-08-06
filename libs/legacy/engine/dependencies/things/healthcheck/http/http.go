package httpHealthcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ongruent/gruent/libs/legacy/engine/dependencies/things"
)

func HttpServer(attrs *things.HttpServer) error {
	// Initialize HTTP client
	client := http.Client{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Send an HTTP request
	request, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s://%s:%d", attrs.Protocol, attrs.Host, attrs.Port), nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP server returned status code %d", response.StatusCode)
	}

	return nil
}
