package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server Server `yaml:"server"`
}

type Server struct {
	Port string `env:"PORT" yaml:"port" env-default:"8080"`
}

func Read(log *log.Logger, yamlPath string) *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Printf("read config from env err: %v\n", err)
	}

	err = cleanenv.ReadConfig(yamlPath, &cfg)
	if err != nil {
		log.Printf("read config err: %v\n", err)
	}

	return &cfg
}
