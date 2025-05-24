package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lerentis/stunning-fiesta/internal/config"
)

func FetchTemplatesList(templatesEndpoint string) (*TemplatesListResponse, error) {
	resp, err := http.Get(templatesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch templates list: %w", err)
	}
	defer resp.Body.Close()

	var list TemplatesListResponse
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode templates list: %w", err)
	}
	return &list, nil
}

func DownloadTemplate(url, dir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download template %s: %w", url, err)
	}
	defer resp.Body.Close()

	filename := filepath.Base(url)
	outPath := filepath.Join(dir, filename)
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", outPath, err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", outPath, err)
	}
	return outPath, nil
}

func GetTemplates(cfg config.Config) (string, error) {
	list, err := FetchTemplatesList(cfg.Endpoints.Template)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "stunning-fiesta-templates-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	for _, url := range list.Templates {
		if _, err := DownloadTemplate(url, tempDir); err != nil {
			return "", err
		}
	}

	return tempDir, nil
}
