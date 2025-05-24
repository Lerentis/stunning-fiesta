package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableProjects(cfg config.Config) (*ProjectList, error) {
	if cfg.Endpoints.Projects == "" {
		return nil, fmt.Errorf("Projects Endpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.DBMS)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Projects: %w", err)
	}
	defer resp.Body.Close()

	var list ProjectList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode Project list: %w", err)
	}

	return &list, nil
}
