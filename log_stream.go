package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ravelAPI "github.com/valyentdev/ravel/api"
)

// LogStreamOptions defines the parameters for streaming logs
type LogStreamOptions struct {
	FleetID   string
	MachineID string
	Follow    bool
	Namespace string
}

// LogStream represents an active log streaming session
type LogStream struct {
	reader   *bufio.Reader
	response *http.Response
	err      error
}

// Close cleanly terminates the log stream
func (ls *LogStream) Close() error {
	if ls.response != nil && ls.response.Body != nil {
		return ls.response.Body.Close()
	}
	return nil
}

// Next returns the next log entry from the stream.
// Returns false when the stream ends or an error occurs.
// Use Err() to check for errors after Next() returns false.
func (ls *LogStream) Next() (ravelAPI.LogEntry, bool) {
	line, err := ls.reader.ReadBytes('\n')
	if err != nil {
		if err != io.EOF {
			ls.err = fmt.Errorf("failed to read log line: %w", err)
		}
		return ravelAPI.LogEntry{}, false
	}

	var entry ravelAPI.LogEntry
	if err := json.Unmarshal(line, &entry); err != nil {
		ls.err = fmt.Errorf("failed to unmarshal log entry: %w", err)
		return ravelAPI.LogEntry{}, false
	}

	return entry, true
}

// Err returns any error that occurred during streaming
func (ls *LogStream) Err() error {
	return ls.err
}

// StreamLogs initiates a streaming connection to receive logs in real-time
func (client *Client) StreamLogs(ctx context.Context, opts LogStreamOptions) (*LogStream, error) {
	// Construct the path with query parameters
	path := fmt.Sprintf("/v1/fleets/%s/machines/%s/logs", opts.FleetID, opts.MachineID)
	if opts.Follow {
		path += "?follow=true"
	}
	if opts.Namespace != "" {
		if opts.Follow {
			path += "&"
		} else {
			path += "?"
		}
		path += fmt.Sprintf("namespace=%s", opts.Namespace)
	}

	// Create the HTTP request
	url := client.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header if available
	if client.bearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+client.bearerToken)
	}

	// Set appropriate headers for streaming
	req.Header.Set("Accept", "application/x-ndjson")

	// Send the request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create buffered reader for the response body
	reader := bufio.NewReaderSize(resp.Body, 64*1024) // 64KB buffer

	return &LogStream{
		reader:   reader,
		response: resp,
	}, nil
}
