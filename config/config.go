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

	// LDAP configuration
	LDAPServer           string        `json:"ldap_server" yaml:"ldap_server" env:"LDAP_SERVER"`
	LDAPBaseDN           string        `json:"ldap_base_dn" yaml:"ldap_base_dn" env:"LDAP_BASE_DN"`
	LDAPBindDN           string        `json:"ldap_bind_dn" yaml:"ldap_bind_dn" env:"LDAP_BIND_DN"`
	LDAPBindPassword     string        `json:"ldap_bind_password" yaml:"ldap_bind_password" env:"LDAP_BIND_PASSWORD"`
	LDAPUserSearchFilter string        `json:"ldap_user_search_filter" yaml:"ldap_user_search_filter" env:"LDAP_USER_SEARCH_FILTER"`
	LDAPUserSearchAttr   string        `json:"ldap_user_search_attr" yaml:"ldap_user_search_attr" env:"LDAP_USER_SEARCH_ATTR"`
	LDAPUseTLS           bool          `json:"ldap_use_tls" yaml:"ldap_use_tls" env:"LDAP_USE_TLS"`
	LDAPInsecureSkipTLS  bool          `json:"ldap_insecure_skip_tls" yaml:"ldap_insecure_skip_tls" env:"LDAP_INSECURE_SKIP_TLS"`
	LDAPTimeout          time.Duration `json:"ldap_timeout" yaml:"ldap_timeout" env:"LDAP_TIMEOUT"`
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

		// LDAP configuration
		LDAPServer:           getEnvOrDefault("LDAP_SERVER", defaults.LDAPServer),
		LDAPBaseDN:           getEnvOrDefault("LDAP_BASE_DN", defaults.LDAPBaseDN),
		LDAPBindDN:           getEnvOrDefault("LDAP_BIND_DN", defaults.LDAPBindDN),
		LDAPBindPassword:     getEnvOrDefault("LDAP_BIND_PASSWORD", defaults.LDAPBindPassword),
		LDAPUserSearchFilter: getEnvOrDefault("LDAP_USER_SEARCH_FILTER", defaults.LDAPUserSearchFilter),
		LDAPUserSearchAttr:   getEnvOrDefault("LDAP_USER_SEARCH_ATTR", defaults.LDAPUserSearchAttr),
		LDAPUseTLS:           getEnvBoolOrDefault("LDAP_USE_TLS", defaults.LDAPUseTLS),
		LDAPInsecureSkipTLS:  getEnvBoolOrDefault("LDAP_INSECURE_SKIP_TLS", defaults.LDAPInsecureSkipTLS),
		LDAPTimeout:          getEnvDurationOrDefault("LDAP_TIMEOUT", defaults.LDAPTimeout),
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

		// LDAP defaults
		LDAPServer:           "ldap://localhost:389",
		LDAPBaseDN:           "dc=example,dc=com",
		LDAPBindDN:           "",
		LDAPBindPassword:     "",
		LDAPUserSearchFilter: "(uid=%s)",
		LDAPUserSearchAttr:   "uid",
		LDAPUseTLS:           false,
		LDAPInsecureSkipTLS:  false,
		LDAPTimeout:          10 * time.Second,
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

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
