package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const basicConfigPath = "./configs/app.yaml"

var configPath string

type Config struct {
	Server HTTPServer `yaml:"server"`
}

func NewConfig() *Config {
	var cfg Config

	setupConfigPath()
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = basicConfigPath
		}
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		msg := fmt.Sprintf(`
Failed to read config. Error: %s
Use --config flag to setup config path.
Example usage: go run ./cmd/api/main.go --config="./configs/app.yaml"
	`, err.Error())
		panic(msg)
	}

	return &cfg
}

func setupConfigPath() {
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()
}
