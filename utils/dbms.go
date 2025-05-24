package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableDBMS(cfg config.Config) (*DBMSList, error) {
	if cfg.Endpoints.DBMS == "" {
		return nil, fmt.Errorf("DBMS Endpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.DBMS)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DBMS's: %w", err)
	}
	defer resp.Body.Close()

	var list DBMSList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode DBMS list: %w", err)
	}

	return &list, nil
}
