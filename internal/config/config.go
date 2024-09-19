package config

import (
	"log"

	"github.com/shuryak/sberhack/pkg/configman"
)

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Port string `env:"PORT" yaml:"port" default:"8080"`
}

func Read(log *log.Logger, yamlPath string) *Config {
	var cfg Config

	err := configman.Collect(&cfg, configman.WithYAMLFile(yamlPath), configman.WithEnv())
	if err != nil {
		log.Fatalf("read config from env, err: %v\n", err)
	}

	return &cfg
}
