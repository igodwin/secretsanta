package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	configInstance *Config
	once           sync.Once
	Paths          = []string{"/etc/secretsanta/", "$HOME/.secretsanta", ".", getBinaryDir()}
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

	// Try to find config file with name "secretsanta.config" (no extension)
	viper.SetConfigName("secretsanta.config")
	viper.SetConfigType("toml")

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

	return config
}
