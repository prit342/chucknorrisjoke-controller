package chuckclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

// APIClient is a client for the Chuck Norris API
type APIClient struct{}

// APIResponse is the response from the Chuck Norris API
type APIResponse struct {
	Categories []string `json:"categories"`
	CreatedAt  string   `json:"created_at"`
	IconURL    string   `json:"icon_url"`
	ID         string   `json:"id"`
	UpdatedAt  string   `json:"updated_at"`
	URL        string   `json:"url"`
	Value      string   `json:"value"`
}

// APIURL is the URL of the Chuck Norris API
const APIURL = "https://api.chucknorris.io/jokes/random?category="

// Categories is a list of categories that is supported by the Chuck Norris API
var categories = []string{
	"animal",
	"career",
	"celebrity",
	"dev",
	"explicit",
	"fashion",
	"food",
	"history",
	"money",
	"movie",
	"music",
	"political",
	"religion",
	"science",
	"sport",
	"travel",
}

// JokeGetter is an interface for getting jokes
type JokeGetter interface {
	GetJoke(ctx context.Context, category string) (string, error)
}

// GetJoke returns a joke from the Chuck Norris API
func (c *APIClient) GetJoke(ctx context.Context, category string) (string, error) {

	if !slices.Contains(categories, category) {
		return "", fmt.Errorf("category %s is not supported", category)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, APIURL+category, nil)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}

	defer func() { _ = resp.Body.Close() }()

	var response APIResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("could not decode response: %w", err)
	}

	return response.Value, nil
}
