package cmd

import "encoding/json"

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func formatJSONOutput(data []map[string]interface{}) (string, error) {
	output, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return "", err
	}

	return string(output), nil
}
