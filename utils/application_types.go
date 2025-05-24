package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableApplicationTypes(cfg config.Config) (*ApplicationTypesList, error) {
	if cfg.Endpoints.ApplicationTypes == "" {
		return nil, fmt.Errorf("ApplicationTypesEndpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.ApplicationTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Application Types: %w", err)
	}
	defer resp.Body.Close()

	var list ApplicationTypesList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode Application Types list: %w", err)
	}

	return &list, nil
}
