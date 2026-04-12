package resources

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var HTTPClient = &http.Client{
	Timeout: 10 * time.Second,
}

var (
	cache  = make(map[string][]byte)
	cachMu sync.RWMutex
)

func HttpGet(url string) ([]byte, error) {
	cachMu.RLock()
	if body, ok := cache[url]; ok {
		cachMu.RUnlock()
		return body, nil
	}
	cachMu.RUnlock()

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

	cachMu.Lock()
	cache[url] = body
	cachMu.Unlock()

	return body, nil
}
