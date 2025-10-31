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

// ResetConfig resets the singleton config instance (useful for testing)
func ResetConfig() {
	configInstance = nil
	once = sync.Once{}
	viper.Reset()
}

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
	ContentType string `mapstructure:"content_type"`
}

type NotifierConfig struct {
	ServiceAddr string `mapstructure:"service_addr"`
	ArchiveEmail string `mapstructure:"archive_email"`
	APIKey      string `mapstructure:"api_key"`
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
	viper.SetDefault("smtp.content_type", "text/plain")
	viper.SetDefault("notifier.service_addr", "")
	viper.SetDefault("notifier.archive_email", "")
	viper.SetDefault("notifier.api_key", "")

	viper.AutomaticEnv()

	// Try to find config file - look for config.yaml first, then secretsanta.config
	configFound := false
	configPath := ""
	var configErr error

	// Try config.yaml first
	for _, searchPath := range Paths {
		potentialPath := filepath.Join(searchPath, "config.yaml")
		if _, err := os.Stat(potentialPath); err == nil {
			viper.SetConfigFile(potentialPath)
			viper.SetConfigType("yaml")
			if err := viper.ReadInConfig(); err == nil {
				configPath = potentialPath
				configFound = true
				break
			} else {
				configErr = err
			}
		}
	}

	// Fall back to secretsanta.config if config.yaml not found
	if !configFound {
		for _, searchPath := range Paths {
			potentialPath := filepath.Join(searchPath, "secretsanta.config")
			if _, err := os.Stat(potentialPath); err == nil {
				viper.SetConfigFile(potentialPath)
				viper.SetConfigType("toml")
				if err := viper.ReadInConfig(); err == nil {
					configPath = potentialPath
					configFound = true
					break
				} else {
					configErr = err
				}
			}
		}
	}

	if configFound {
		log.Printf("Loaded config from: %s", configPath)
	} else {
		if configErr != nil {
			log.Printf("Error reading config file: %v", configErr)
		}
		log.Println("No config file found, using defaults (stdout notifications only)")
	}

	config := &Config{}
	err := viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}

	// Log the loaded configuration
	logConfig(config, configPath)

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
			"content_type": cfg.SMTP.ContentType,
		},
		"notifier": map[string]interface{}{
			"service_addr":  cfg.Notifier.ServiceAddr,
			"archive_email": cfg.Notifier.ArchiveEmail,
			"api_key":       redact(cfg.Notifier.APIKey),
		},
	}

	jsonBytes, err := json.MarshalIndent(redactedConfig, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal config for logging: %v", err)
		return
	}

	log.Printf("Configuration:\n%s", string(jsonBytes))
}
