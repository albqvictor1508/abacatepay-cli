package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.abacatepay.com/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func ValidateAPIKey(apiKey string) (bool, error) {
	client := NewClient(apiKey)

	req, err := http.NewRequest("GET", baseURL+"/store/get", nil)
	if err != nil {
		return false, nil
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	reply, err := client.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer reply.Body.Close()

	if reply.StatusCode == http.StatusOK {
		return true, nil
	}

	if reply.StatusCode == http.StatusUnauthorized {
		return false, nil
	}

	var error map[string]interface{}
	if err := json.NewDecoder(reply.Body).Decode(&error); err != nil {
		return false, fmt.Errorf("error to validate: status %d", reply.StatusCode)
	}

	return false, fmt.Errorf("error to validate: %v", error)
}
