package api

import (
	"bytes"
	"encoding/json"
	"errors"
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
