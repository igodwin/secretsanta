package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var configInstance *Config
var once sync.Once

type Config struct {
	SMTP struct {
		Host        string `mapstructure:"host"`
		Port        string `mapstructure:"port"`
		Identity    string `mapstructure:"identity"`
		Username    string `mapstructure:"username"`
		Password    string `mapstructure:"password"`
		FromAddress string `mapstructure:"from_address"`
		FromName    string `mapstructure:"from_name"`
	} `mapstructure:"smtp"`
	OtherSetting string `mapstructure:"otherSetting"`
}

func (config *Config) SMTPIsConfigured() bool {
	return config.SMTP.Host != "" && config.SMTP.Port != "" && config.SMTP.Username != "" && config.SMTP.Password != "" && config.SMTP.FromAddress != ""
}

func GetConfig() *Config {
	once.Do(func() {
		configInstance = loadConfig()
	})
	return configInstance
}

func loadConfig() *Config {
	viper.SetConfigFile("secretsanta.config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/secretsanta/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.secretsanta") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return config
}
