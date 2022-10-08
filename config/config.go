package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	API API
	DB  DB
}

type API struct {
	BaseUrl string
	Port    string
}

type DB struct {
	Host     string
	Username string
	Password string
	BaseUrl  string
	Database string
}

func NewConfig() *Config {
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return &Config{
		API: API{
			BaseUrl: viper.GetString("API.BaseURL"),
			Port:    viper.GetString("API.Port"),
		},
		DB: DB{
			Host:     viper.GetString("DB.Host"),
			Username: viper.GetString("DB.Username"),
			Password: viper.GetString("DB.Password"),
			BaseUrl:  viper.GetString("DB.BaseURL"),
			Database: viper.GetString("DB.Database"),
		},
	}
}
