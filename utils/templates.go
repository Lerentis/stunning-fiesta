package utils

import (
	"encoding/json"
	"fmt"
	"html/template"
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

func GetNamespaceTemplates(cfg config.Config) (string, error) {
	list, err := FetchTemplatesList(cfg.Endpoints.Template)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "stunning-fiesta-namespace-templates-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	baseURL := cfg.Endpoints.Template
	if idx := len(baseURL) - len(filepath.Base(baseURL)); idx > 0 {
		baseURL = baseURL[:idx]
	}

	for _, filename := range list.NamespaceTemplates {
		url := baseURL + "namespace-templates/" + filename
		if _, err := DownloadTemplate(url, tempDir); err != nil {
			return "", err
		}
	}

	return tempDir, nil
}

func GetHelmTemplates(cfg config.Config) (string, error) {
	list, err := FetchTemplatesList(cfg.Endpoints.Template)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "stunning-fiesta-helm-templates-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	baseURL := cfg.Endpoints.Template
	if idx := len(baseURL) - len(filepath.Base(baseURL)); idx > 0 {
		baseURL = baseURL[:idx]
	}

	for _, filename := range list.HelmTemplates {
		url := baseURL + "helm-templates/" + filename
		if _, err := DownloadTemplate(url, tempDir); err != nil {
			return "", err
		}
	}

	return tempDir, nil
}

func GetApplicationTemplates(cfg config.Config) (string, error) {
	list, err := FetchTemplatesList(cfg.Endpoints.Template)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "stunning-fiesta-application-templates-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	baseURL := cfg.Endpoints.Template
	if idx := len(baseURL) - len(filepath.Base(baseURL)); idx > 0 {
		baseURL = baseURL[:idx]
	}

	for _, filename := range list.ApplicationTemplates {
		url := baseURL + "application-templates/" + filename
		if _, err := DownloadTemplate(url, tempDir); err != nil {
			return "", err
		}
	}

	return tempDir, nil
}

func GetInfrastructureTemplates(cfg config.Config) (string, error) {
	list, err := FetchTemplatesList(cfg.Endpoints.Template)
	if err != nil {
		return "", err
	}

	tempDir, err := os.MkdirTemp("", "stunning-fiesta-infrastructure-templates-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	baseURL := cfg.Endpoints.Template
	if idx := len(baseURL) - len(filepath.Base(baseURL)); idx > 0 {
		baseURL = baseURL[:idx]
	}

	for _, filename := range list.InfrastructureTemplates {
		url := baseURL + "infrastructure-templates/" + filename
		if _, err := DownloadTemplate(url, tempDir); err != nil {
			return "", err
		}
	}

	return tempDir, nil
}

func RenderTemplatesDir(srcDir, dstDir string, vars map[string]interface{}) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dstDir, relPath)

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}
		outFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}
		defer outFile.Close()

		if err := tmpl.Execute(outFile, vars); err != nil {
			return fmt.Errorf("failed to render template %s: %w", path, err)
		}
		return nil
	})
}
