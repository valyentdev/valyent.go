package valyent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/valyentdev/ravel/api"
)

type Deployment struct {
	ID     string `json:"id"`
	Origin string `json:"origin"`
	Status string `json:"status"`
}

type CreateDeploymentPayload struct {
	Machine api.CreateMachinePayload
}

func (client *Client) CreateDeployment(namespace, fleetID string,
	payload CreateDeploymentPayload,
	tarball io.ReadCloser,
) (*Deployment, error) {
	url := client.baseURL + "/organizations/" + namespace + "/applications/" + fleetID + "/deployments"

	var req *http.Request
	var err error

	if tarball != nil {
		// Create a multipart form
		form := bytes.NewBuffer(nil)
		writer := multipart.NewWriter(form)

		// Add machine data to the form
		machine, err := json.Marshal(payload.Machine)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal machine data: %v", err)
		}

		if err := writer.WriteField("machine", string(machine)); err != nil {
			return nil, fmt.Errorf("failed to write machine data to request: %v", err)
		}

		// Add tarball to the form
		part, err := writer.CreateFormFile("tarball", "tarball")
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(part, tarball)
		if err != nil {
			return nil, err
		}

		// Close the writer
		err = writer.Close()
		if err != nil {
			return nil, err
		}

		// Create a new HTTP request with multipart data
		req, err = http.NewRequest("POST", url, form)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", writer.FormDataContentType())
	} else {
		// Use JSON payload if tarball is nil
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON payload: %v", err)
		}

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
	}

	// Common headers
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+client.bearerToken)

	// Send the HTTP request
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 response
	if resp.StatusCode != 200 {
		if resp.Body != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("HTTP request failed with status code %d with body: %s", resp.StatusCode, bodyBytes)
		}
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	// Decode deployment as JSON
	deployment := &Deployment{}
	if err := json.NewDecoder(resp.Body).Decode(deployment); err != nil {
		return nil, err
	}

	return deployment, nil
}
