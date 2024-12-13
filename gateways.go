package api

import (
	"fmt"
	stdHTTP "net/http"

	"github.com/valyentdev/ravel/api"
)

func (client *Client) CreateGateway(payload api.CreateGatewayPayload) (*api.Gateway, error) {
	gateway := &api.Gateway{}
	err := client.PerformRequest("POST", "/v1/gateways", payload, &gateway)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway from the api: %v", err)
	}
	return gateway, nil
}

func (client *Client) GetGateways() ([]api.Gateway, error) {
	gateways := []api.Gateway{}
	err := client.PerformRequest("GET", "/v1/gateways", nil, &gateways)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve gateways from the api: %v", err)
	}
	return gateways, nil
}

func (client *Client) DeleteGateway(gatewayID string) error {
	err := client.PerformRequest(
		stdHTTP.MethodDelete,
		"/v1/gateways/"+gatewayID,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to delete gateway: %v", err)
	}
	return err
}
