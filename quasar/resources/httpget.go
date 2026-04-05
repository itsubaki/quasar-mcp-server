package resources

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var HTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

func HttpGet(url string) ([]byte, error) {
	resp, err := HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read all: %w", err)
	}

	return body, nil
}
