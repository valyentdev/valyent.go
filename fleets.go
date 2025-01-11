package valyent

import (
	"fmt"
	stdHTTP "net/http"

	"github.com/valyentdev/ravel/api"
)

func (client *Client) CreateFleet(opts api.CreateFleetPayload) (*api.Fleet, error) {
	fleet := &api.Fleet{}
	err := client.PerformRequest(
		"POST",
		"/v1/fleets",
		opts,
		fleet,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create fleet: %v", err)
	}

	return fleet, nil
}

func (client *Client) GetFleets() (fleets []api.Fleet, err error) {
	// Fetch existing fleets matching the user's namespace.
	fleets = []api.Fleet{}
	err = client.PerformRequest(stdHTTP.MethodGet, "/v1/fleets", nil, &fleets)
	if err != nil {
		return fleets, fmt.Errorf("failed to retrieve fleets: %v", err)
	}

	return
}

func (client *Client) DeleteFleet(fleetID string) error {
	err := client.PerformRequest(stdHTTP.MethodDelete, "/v1/fleets/"+fleetID, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete fleet: %v", err)
	}

	return nil
}
