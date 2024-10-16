package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	configInstance *Config
	once           sync.Once
	Paths          = []string{"/etc/secretsanta/", "$HOME/.secretsanta", "."}
)

type Config struct {
	SMTP SMTPConfig `mapstructure:"smtp"`
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

func GetConfig() *Config {
	once.Do(func() {
		configInstance = loadConfig()
	})
	return configInstance
}

func loadConfig() *Config {
	viper.SetConfigName("secretsanta.config")
	viper.SetConfigType("toml")

	for _, path := range Paths {
		viper.AddConfigPath(path)
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}

	return config
}
