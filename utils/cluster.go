package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func ListAvailableClusters(cfg config.Config) (*ClusterList, error) {
	if cfg.Endpoints.ClusterInfo == "" {
		return nil, fmt.Errorf("ClusterInfoEndpoint is not set in config")
	}

	resp, err := http.Get(cfg.Endpoints.ClusterInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clusters: %w", err)
	}
	defer resp.Body.Close()

	var list ClusterList
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode clusters list: %w", err)
	}

	return &list, nil
}
