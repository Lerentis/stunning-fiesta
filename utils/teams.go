package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableTeams(cfg config.Config) (*TeamList, error) {
	if cfg.Endpoints.Teams == "" {
		return nil, fmt.Errorf("Teams Endpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.DBMS)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Teams: %w", err)
	}
	defer resp.Body.Close()

	var list TeamList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode Teams list: %w", err)
	}

	return &list, nil
}
