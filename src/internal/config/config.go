package config

import (
	"github.com/caarlos0/env/v9"
)

// type MonitoringConfig struct {
// 	TraceDestination   string `env:"TRACE_DESTINATION"`
// 	MetricsDestination string `env:"METRICS_DESTINATION"`
// 	LogFileName        string `env:"LOG_FILE_NAME"`
// }

type GenericConfig struct {
	// DatabaseDSN string `env:"DATABASE_DSN"`
	ServicePort uint16 `env:"SERVICE_PORT"`
	ServiceBind string `env:"SERVICE_BIND"`
	Environment string `env:"ENVIRONMENT"`
}

type Config struct {
	// MonitoringConfig
	GenericConfig
}

func InitConfig() (cfg *Config, err error) {
	cfg = &Config{}
	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return
}
