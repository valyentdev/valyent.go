package api

import (
	"fmt"
	stdHTTP "net/http"

	"github.com/valyentdev/ravel/api"
)

func (client *Client) CreateMachine(fleetID string, opts api.CreateMachinePayload) (machine *api.Machine, err error) {
	err = client.PerformRequest("POST", fmt.Sprintf("/v1/fleets/%s/machines", fleetID), opts, machine)
	if err != nil {
		return nil, err
	}

	return machine, nil
}

func (client *Client) GetMachines(fleetID string) ([]api.Machine, error) {
	machines := []api.Machine{}
	err := client.PerformRequest("GET", fmt.Sprintf("/v1/fleets/%s/machines", fleetID), nil, &machines)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve machines from the api: %v", err)
	}
	return machines, nil
}

func (client *Client) GetMachineEvents(fleetID, machineID string) ([]api.MachineEvent, error) {
	events := []api.MachineEvent{}
	err := client.PerformRequest("GET", fmt.Sprintf("/v1/fleets/%s/machines/%s/events", fleetID, machineID), nil, &events)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve events from the api: %v", err)
	}
	return events, nil
}

func (client *Client) DeleteMachine(fleetID, machineID string, force bool) error {
	// Initialize response data struct.
	resp := map[string]any{}

	// Compute request's path.
	path := fmt.Sprintf("/v1/fleets/%s/machines/%s?force=%t", fleetID, machineID, force)

	// Actually call the API.
	err := client.PerformRequest(
		stdHTTP.MethodDelete,
		path,
		nil, &resp,
	)
	if err != nil {
		return fmt.Errorf("failed to delete machine: %v", resp["detail"])
	}

	return nil
}

func (client *Client) GetLogs(fleetID, machineID string) ([]api.LogEntry, error) {
	// Initialize a list of log entries to later be filled by the actual values from the API.
	logEntries := []api.LogEntry{}

	// Let's fetch the actual HTTP API for log entries.
	err := client.PerformRequest(
		"GET",
		fmt.Sprintf("/v1/fleets/%s/machines/%s/logs", fleetID, machineID),
		nil, &logEntries,
	)
	if err != nil {
		return nil, err
	}

	return logEntries, nil
}

func (client *Client) StartMachine(fleetID, machineID string) error {
	err := client.PerformRequest("POST", fmt.Sprintf("/fleets/%s/machines/%s/start", fleetID, machineID), nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (client *Client) StopMachine(fleetID, machineID string) error {
	err := client.PerformRequest("POST", fmt.Sprintf("/fleets/%s/machines/%s/stop", fleetID, machineID), nil, nil)
	if err != nil {
		return err
	}

	return nil
}
