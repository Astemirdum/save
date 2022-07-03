package config

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr       string `envconfig:"SRV_ADDR"`
	Port1      int    `envconfig:"SRV_PORT1"`
	Port2      int    `envconfig:"SRV_PORT2"`
	ClientName string `envconfig:"CLIENT_NAME"`
	Key        string `envconfig:"KEY"`
}

// NewConfig reads config from environment.
func NewConfig() *Config {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}
	configBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Configuration:", string(configBytes))
	return &config
}
