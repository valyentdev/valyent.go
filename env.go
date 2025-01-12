package valyent

import (
	"errors"
	"fmt"
	"strings"
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

func (client *Client) SetEnvironmentVariables(
	namespace string,
	fleetId string,
	args []string,
) (bool, error) {
	// data := "{"
	// for i, arg := range args {
	// 	// Check arg follows the pattern <key>=<value>
	// 	if !strings.Contains(arg, "=") {
	// 		return false, errors.New("invalid argument: " + arg)
	// 	}

	// 	splittedArg := strings.Split(arg, "=")
	// 	key := splittedArg[0]
	// 	value := splittedArg[1]

	// 	// Remove single or double quotes to value
	// 	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
	// 		value = strings.TrimPrefix(value, "\"")
	// 		value = strings.TrimSuffix(value, "\"")
	// 	} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
	// 		value = strings.TrimPrefix(value, "'")
	// 		value = strings.TrimSuffix(value, "'")
	// 	}

	// 	data += "\"" + key + "\":\"" + value + "\""
	// 	if i < len(args)-1 {
	// 		data += ","
	// 	}
	// }
	// data += "}"

	// body := bytes.NewBufferString(data)

	// url := RetrieveApiBaseUrl() + "/applications/" + applicationId + "/env"
	// req, err := http.NewRequest("PATCH", url, body)
	// if err != nil {
	// 	return false, err
	// }

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return false, err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != 200 {
	// 	return false, errors.New("failed to set environment variable")
	// }

	payload := map[string]string{}
	for _, arg := range args {
		// Check arg follows the pattern <key>=<value>
		if !strings.Contains(arg, "=") {
			return false, errors.New("invalid argument: " + arg)
		}

		splittedArg := strings.Split(arg, "=")
		key := splittedArg[0]
		value := splittedArg[1]

		// Remove single or double quotes to value
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.TrimPrefix(value, "\"")
			value = strings.TrimSuffix(value, "\"")
		} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") {
			value = strings.TrimPrefix(value, "'")
			value = strings.TrimSuffix(value, "'")
		}

		payload[key] = value
	}

	response := struct {
		Redeploy bool `json:"redeploy"`
	}{}

	client.PerformRequest("PATCH", "/organizations/"+namespace+"/applications/"+fleetId+"/env", payload, response)

	return response.Redeploy, nil
}
