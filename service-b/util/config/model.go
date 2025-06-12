package config

// App config

type App struct {
	Name               string `mapstructure:"name"`
	Host               string `mapstructure:"host"` // Bind address (0.0.0.0 for listening)
	Port               int    `mapstructure:"port"`
	RegisterAddress    string `mapstructure:"register_address"`     // Address for service registration
	HealthCheckAddress string `mapstructure:"health_check_address"` // Address for Consul health checks
}

// Open telemetry config

type OpenTelemetry struct {
	TracerName   string `mapstructure:"tracer_name"`
	CollectorURL string `mapstructure:"collector_url"`
}
