package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableVaults(cfg config.Config) (*VaultList, error) {
	if cfg.Endpoints.Vault == "" {
		return nil, fmt.Errorf("VaultEndpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.ApplicationTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Vault Data: %w", err)
	}
	defer resp.Body.Close()

	var list VaultList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode Vault Data list: %w", err)
	}

	return &list, nil
}
