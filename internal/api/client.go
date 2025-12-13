package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.abacatepay.com/v1"

type Client struct {
	key        string
	httpClient *http.Client
}

type AbacatePayResponse struct {
	Data  map[string]any `json:"data"`
	Error string         `json:"error"`
}

type RequestOptions struct {
	Route  string
	Method string
	Body   *[]byte
}

func NewClient(key string) *Client {
	return &Client{
		key: key,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Request(options RequestOptions) (map[string]any, error) {
	var bodyReader io.Reader
	if options.Body != nil {
		bodyReader = bytes.NewBuffer(*options.Body)
	}

	req, err := http.NewRequest(options.Method, baseURL+options.Route, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.key)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var payload AbacatePayResponse

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if payload.Error != "" {
		return nil, errors.New(payload.Error)
	}

	return payload.Data, nil
}

func ValidateAPIKey(apiKey string) (bool, error) {
	client := NewClient(apiKey)

	req, err := http.NewRequest("GET", baseURL+"/store/get", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return false, nil
	}

	var errorResponse map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		return false, fmt.Errorf("error to validate: status %d", resp.StatusCode)
	}

	return false, fmt.Errorf("error to validate: %v", errorResponse)
}
