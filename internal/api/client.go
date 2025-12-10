package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.abacatepay.com/v1"

type Client struct {
	key     string
	httpClient *http.Client
}

type AbacatePayResponse struct {
	data	any
	error	string
}

func NewClient(key string) *Client {
	return &Client{
		key: key,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func ValidateAPIKey(key string) (bool, error) {
	client := NewClient(key)

	req, err := http.NewRequest("GET", baseURL + "/store/get", nil)

	if err != nil {
		return false, nil
	}

	req.Header.Set("Authorization", "Bearer " + key)
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

	var error AbacatePayResponse

	if err := json.NewDecoder(reply.Body).Decode(&error); err != nil {
		return false, fmt.Errorf("Unknown AbacatePay error, status (%d)", reply.StatusCode)
	}

	return false, fmt.Errorf(error.error)
}
