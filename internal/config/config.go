package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Endpoints struct {
	Template         string `yaml:"template"`
	Update           string `yaml:"update"`
	ClusterInfo      string `yaml:"cluster_info"`
	Applications     string `yaml:"applications"`
	DBMS             string `yaml:"dbms"`
	Projects         string `yaml:"projects"`
	Vault            string `yaml:"vault"`
	Teams            string `yaml:"teams"`
	ApplicationTypes string `yaml:"application_types"`
}

type Config struct {
	Endpoints Endpoints `yaml:"endpoints"`
	GitlabURL string    `yaml:"gitlab_url"`
}

func LoadConfig(configPath string) (*Config, error) {
	err := EnsureConfig(configPath)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func EnsureConfig(configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultCfg := Config{
			Endpoints: Endpoints{
				Template:         "https://lerentis.github.io/stunning-fiesta/templates.json",
				Update:           "https://lerentis.github.io/stunning-fiesta/version.json",
				ClusterInfo:      "https://lerentis.github.io/stunning-fiesta/clusters.json",
				Applications:     "https://lerentis.github.io/stunning-fiesta/applications.json",
				DBMS:             "https://lerentis.github.io/stunning-fiesta/dbms.json",
				Projects:         "https://lerentis.github.io/stunning-fiesta/projects.json",
				Vault:            "https://lerentis.github.io/stunning-fiesta/vault.json",
				Teams:            "https://lerentis.github.io/stunning-fiesta/teams.json",
				ApplicationTypes: "https://lerentis.github.io/stunning-fiesta/applicationTypes.json",
			},
			GitlabURL: "https://gitlab.com/lerentis/stunning-fiesta",
		}
		data, err := yaml.Marshal(&defaultCfg)
		if err != nil {
			return err
		}
		return os.WriteFile(configPath, data, 0644)
	}
	return nil
}
