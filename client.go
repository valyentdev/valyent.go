package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	DEFAULT_VALYENT_API_BASE_URL = "https://console.valyent.cloud"
)

// Client represents a Valyent HTTP API client.
type Client struct {
	baseURL     string
	bearerToken string
}

// NewClient returns an instance of a Valyent HTTP API client.
func NewClient() *Client {
	return &Client{
		baseURL: DEFAULT_VALYENT_API_BASE_URL, // We set the default value for the base URL.
	}
}

// WithBearerToken allows to specify the authorization bearer token to the HTTP requests.
func (client *Client) WithBearerToken(bearerToken string) *Client {
	client.bearerToken = bearerToken
	return client
}

// WithBaseURL allows to customize the API base url. If not specified, the client will make use of the default value.
func (client *Client) WithBaseURL(baseURL string) *Client {
	client.baseURL = baseURL
	return client
}

// PerformRequest sends an HTTP request with an optional JSON payload and unmarshals the JSON response if responseTarget is non-nil.
func (client *Client) PerformRequest(method, path string, payload any, responseTarget any) error {
	// Construct the URL
	url := client.baseURL + path

	// Prepare request body if payload is provided
	var body io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload to JSON: %w", err)
		}
		body = bytes.NewReader(payloadBytes)
	}

	// Create the HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Add authorization header if API key is available
	if client.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+client.bearerToken)
	}

	// Set appropriate headers
	if payload != nil {
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
	}

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read and process the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// If responseTarget is non-nil, unmarshal the JSON response into it
	if responseTarget != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, responseTarget); err != nil {
			return fmt.Errorf("failed to unmarshal JSON response: %w", err)
		}
	}

	// Check for non-2xx HTTP response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status %s", resp.Status)
	}

	return nil
}
