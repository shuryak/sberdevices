package config

import (
	configman2 "github.com/shuryak/sberdevices/internal/pkg/configman"
)

type (
	Config struct {
		Server   Server   `yaml:"server"`
		Postgres Postgres `yaml:"postgres"`
		Auth     Auth     `yaml:"auth"`
		Clients  Clients  `yaml:"clients" required:"true"`
	}

	Server struct {
		Port string `env:"PORT" yaml:"port" default:"8080"`
	}

	Postgres struct {
		User     string `yaml:"user"`
		Password string `env:"POSTGRES_PASSWORD"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"database"`
	}

	Auth struct {
		CodeLength         uint `yaml:"code_length" default:"64"`
		AccessTokenLength  uint `yaml:"access_token_length" default:"128"`
		RefreshTokenLength uint `yaml:"refresh_token_length" default:"128"`
	}

	Clients []Client

	Client struct {
		ID     string `yaml:"id"`
		Secret string `yaml:"secret"`
	}
)

func (c Clients) CheckClient(id, secret string) bool {
	for _, client := range c {
		if client.ID == id {
			if client.Secret == secret {
				return true
			}
			return false
		}
	}
	return false
}

func Read(yamlPath string) (*Config, error) {
	var cfg Config

	err := configman2.Collect(&cfg, configman2.WithYAMLFile(yamlPath), configman2.WithEnv())
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
