package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	configInstance *Config
	once           sync.Once
	Paths          = []string{getBinaryDir(), ".", "$HOME/.secretsanta", "/etc/secretsanta/"}
)

// getBinaryDir returns the directory containing the running binary
func getBinaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

type Config struct {
	SMTP      SMTPConfig      `mapstructure:"smtp"`
	Notifier  NotifierConfig  `mapstructure:"notifier"`
}

type SMTPConfig struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Identity    string `mapstructure:"identity"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	FromAddress string `mapstructure:"from_address"`
	FromName    string `mapstructure:"from_name"`
}

type NotifierConfig struct {
	ServiceAddr    string `mapstructure:"service_addr"`
	ArchiveEmail   string `mapstructure:"archive_email"`
}

func GetConfig() *Config {
	once.Do(func() {
		configInstance = loadConfig()
	})
	return configInstance
}

func loadConfig() *Config {
	// Set defaults for all fields (allows running without config file)
	viper.SetDefault("smtp.host", "")
	viper.SetDefault("smtp.port", "")
	viper.SetDefault("smtp.identity", "")
	viper.SetDefault("smtp.username", "")
	viper.SetDefault("smtp.password", "")
	viper.SetDefault("smtp.from_address", "")
	viper.SetDefault("smtp.from_name", "Secret Santa")
	viper.SetDefault("notifier.service_addr", "")
	viper.SetDefault("notifier.archive_email", "")

	viper.AutomaticEnv()

	// Try to find config file named "config.yaml"
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	for _, path := range Paths {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig()
	if err != nil {
		// Config file not found is acceptable - use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// It's a different error (e.g., parse error)
			log.Fatalf("error reading config file: %v", err)
		}
		log.Println("No config file found, using defaults (stdout notifications only)")
	} else {
		log.Printf("Loaded config from: %s", viper.ConfigFileUsed())
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}

	// Log the loaded configuration
	logConfig(config, viper.ConfigFileUsed())

	return config
}

// redact returns a redacted version of the string for logging
func redact(s string) string {
	if s == "" {
		return ""
	}
	return "***REDACTED***"
}

// logConfig logs the configuration in a structured format with sensitive fields redacted
func logConfig(cfg *Config, configFile string) {
	// Create a redacted version of the config for logging
	redactedConfig := map[string]interface{}{
		"config_file": configFile,
		"smtp": map[string]interface{}{
			"host":         cfg.SMTP.Host,
			"port":         cfg.SMTP.Port,
			"username":     cfg.SMTP.Username,
			"password":     redact(cfg.SMTP.Password),
			"from_address": cfg.SMTP.FromAddress,
			"from_name":    cfg.SMTP.FromName,
			"identity":     cfg.SMTP.Identity,
		},
		"notifier": map[string]interface{}{
			"service_addr":  cfg.Notifier.ServiceAddr,
			"archive_email": cfg.Notifier.ArchiveEmail,
		},
	}

	jsonBytes, err := json.MarshalIndent(redactedConfig, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal config for logging: %v", err)
		return
	}

	log.Printf("Configuration:\n%s", string(jsonBytes))
}
