package ghapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type EventRaw struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload json.RawMessage `json:"payload"`
}

type Client interface {
	UserEvents(ctx context.Context, username string) ([]EventRaw, error)
}

type httpClient struct {
	hc *http.Client
}

func NewClient() Client {
	return &httpClient{
		hc: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *httpClient) UserEvents(ctx context.Context, username string) ([]EventRaw, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// هدرهای لازم برای GitHub API
	req.Header.Set("User-Agent", "github-activity-cli")
	req.Header.Set("Accept", "application/vnd.github+json")

	// اختیاری: توکن برای بالا بردن rate limit
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
		// ok
	case http.StatusNotFound:
		return nil, fmt.Errorf("user '%s' not found (404)", username)
	case http.StatusForbidden:
		return nil, errors.New("API rate limit exceeded or access forbidden (403). Set GITHUB_TOKEN to increase limits")
	default:
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var events []EventRaw
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return events, nil
}
