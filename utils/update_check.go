package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

func UpdateCheck(updateEndpoint string, version string) bool {
	resp, err := http.Get(updateEndpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var data struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return false
	}

	return strings.TrimSpace(version) == strings.TrimSpace(data.Version)
}
