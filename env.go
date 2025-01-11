package valyent

import (
	"fmt"
)

func (client *Client) GetEnvironmentVariables(namespace, fleetID string) (map[string]string, error) {
	response := &struct {
		Env map[string]string `json:"env"`
	}{}

	path := fmt.Sprintf("/organizations/%s/applications/%s/env", namespace, fleetID)
	err := client.PerformRequest("GET", path, nil, response)
	if err != nil {
		return nil, fmt.Errorf("failed to get environment variables: %v", err)
	}

	return response.Env, nil
}
