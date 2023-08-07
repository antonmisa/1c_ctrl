package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v2"
)

// Config -.
type Config struct {
	App   `yaml:"app"`
	Cache `yaml:"cache"`
	HTTP  `yaml:"http"`
	Log   `yaml:"logger"`
}

// App -.
type App struct {
	PathToRAC string `env-required:"true" yaml:"path_to_rac" env:"PATH_TO_RAC"`
	PathTo1C  string `env-required:"true" yaml:"path_to_1c" env:"PATH_TO_1C"`

	LockCode string `env-required:"true" yaml:"lock_code"`
}

// Cache -.
type Cache struct {
	TTL time.Duration `env-required:"true" yaml:"ttl"`
}

// HTTP -.
type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

// Log -.
type Log struct {
	Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	Path  string `env-required:"true" yaml:"path"`
}

func New() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yml"
	}

	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func Prepare() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yml"
	}

	if _, err := os.Stat(configPath); err == nil {
		return os.ErrExist
	}

	cfg := &Config{
		App{
			PathToRAC: "path to rac file",
			PathTo1C:  "path to 1c executable client",
			LockCode:  "12345",
		},
		Cache{
			TTL: 60,
		},
		HTTP{
			Port: "8080",
		},
		Log{
			Level: "debug",
			Path:  "log.log",
		},
	}

	yamlData, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("error while marshaling config: %w", err)
	}

	err = os.WriteFile(configPath, yamlData, 0644)
	if err != nil {
		return fmt.Errorf("error while creating config.yml: %w", err)
	}

	return nil
}
