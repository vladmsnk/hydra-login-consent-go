package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HydraAdminURL string `json:"hydra_admin_url" yaml:"hydra_admin_url" env:"HYDRA_ADMIN_URL"`
	HydraUsername string `json:"hydra_username" yaml:"hydra_username" env:"HYDRA_USERNAME"`
	HydraPassword string `json:"hydra_password" yaml:"hydra_password" env:"HYDRA_PASSWORD"`

	Port int    `json:"port" yaml:"port" env:"PORT"`
	Host string `json:"host" yaml:"host" env:"HOST"`

	LoginURL   string `json:"login_url" yaml:"login_url" env:"LOGIN_URL"`
	ConsentURL string `json:"consent_url" yaml:"consent_url" env:"CONSENT_URL"`

	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout" env:"READ_TIMEOUT"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout" env:"IDLE_TIMEOUT"`

	ShutdownTimeout time.Duration `json:"shutdown_timeout" yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
}

func NewConfig() *Config {
	defaults := GetDefaultConfig()

	config := &Config{
		HydraAdminURL:   getEnvOrDefault("HYDRA_ADMIN_URL", defaults.HydraAdminURL),
		HydraUsername:   getEnvOrDefault("HYDRA_USERNAME", defaults.HydraUsername),
		HydraPassword:   getEnvOrDefault("HYDRA_PASSWORD", defaults.HydraPassword),
		Host:            getEnvOrDefault("HOST", defaults.Host),
		LoginURL:        getEnvOrDefault("LOGIN_URL", defaults.LoginURL),
		ConsentURL:      getEnvOrDefault("CONSENT_URL", defaults.ConsentURL),
		ReadTimeout:     getEnvDurationOrDefault("READ_TIMEOUT", defaults.ReadTimeout),
		WriteTimeout:    getEnvDurationOrDefault("WRITE_TIMEOUT", defaults.WriteTimeout),
		IdleTimeout:     getEnvDurationOrDefault("IDLE_TIMEOUT", defaults.IdleTimeout),
		ShutdownTimeout: getEnvDurationOrDefault("SHUTDOWN_TIMEOUT", defaults.ShutdownTimeout),
	}

	if portStr := os.Getenv("PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil {
			config.Port = port
		} else {
			config.Port = defaults.Port
		}
	} else {
		config.Port = defaults.Port
	}

	return config
}

func GetDefaultConfig() *Config {
	return &Config{
		HydraAdminURL:   "https://oauthidm.ru",
		HydraUsername:   "adminuser",
		HydraPassword:   "1234",
		Port:            3000,
		Host:            "127.0.0.1",
		LoginURL:        "/login",
		ConsentURL:      "/consent",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
