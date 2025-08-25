package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
	"github.com/jstaud/edgelink-go/pkg/models"
)

type Broker struct {
	MQTTURL   string `yaml:"mqtt_url"`
	BaseTopic string `yaml:"base_topic"`
	ClientID  string `yaml:"client_id"`
}

type Config struct {
	Broker       Broker                `yaml:"broker"`
	HTTPAddr     string                `yaml:"http_addr"`
	Devices      []models.DeviceSpec   `yaml:"devices"`
	ShutdownWait time.Duration         `yaml:"shutdown_wait"`
}

func Load(path string) (Config, error) {
	// Read the YAML file
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	
	// Parse YAML into our Config struct
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return Config{}, err
	}
	
	return c, nil
}