package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableProjects(cfg config.Config) (*ProjectList, error) {
	if cfg.Endpoints.Projects == "" {
		return nil, fmt.Errorf("Projects Endpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.Projects)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Projects: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	/*if err != nil {
		return nil, fmt.Errorf("failed to read Projects response body: %w", err)
	}
	fmt.Printf("DEBUG: Projects API response body: %s\n", string(bodyBytes))*/

	var list ProjectList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, fmt.Errorf("failed to decode Project list: %w", err)
	}

	return &list, nil
}
