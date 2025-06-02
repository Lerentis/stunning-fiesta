package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableTeams(cfg config.Config) (*TeamList, error) {
	if cfg.Endpoints.Teams == "" {
		return nil, fmt.Errorf("Teams Endpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.Teams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Teams: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	/*if err != nil {
		return nil, fmt.Errorf("failed to read Teams response body: %w", err)
	}
	fmt.Printf("DEBUG: Teams API response body: %s\n", string(bodyBytes))*/

	var list TeamList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, fmt.Errorf("failed to decode Teams list: %w", err)
	}

	return &list, nil
}
