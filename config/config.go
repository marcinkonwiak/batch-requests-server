package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port         int      `mapstructure:"port"`
	BaseUrl      string   `mapstructure:"base_url"`
	AllowedPaths []string `mapstructure:"allowed_paths"`
}

func NewConfig() *Config {
	config := Config{}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

func LoadConfig() (config *Config) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("port", 1323)
	viper.SetDefault("allowed_paths", []string{})
	viper.SetDefault("base_url", "http://localhost:8080")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	err = viper.SafeWriteConfig()
	if err != nil {
		fmt.Println(err)
	}

	return
}
