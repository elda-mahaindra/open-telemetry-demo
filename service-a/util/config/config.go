package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App           App           `mapstructure:"app"`
	ServiceB      ServiceB      `mapstructure:"service_b"`
	OpenTelemetry OpenTelemetry `mapstructure:"open_telemetry"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to read configuration file: %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal configuration: %s", err)
	}

	return
}
