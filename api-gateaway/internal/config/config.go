package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Environment    string   `mapstructure:"ENVIRONMENT"`
	ServerAddress  string   `mapstructure:"HTTP_SERVER_ADDRESS"`
	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`
}

var conf Configuration

func LoadConfig(path string) (config Configuration, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigFile("gateway.env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return conf, err
	}
	err = viper.Unmarshal(&conf)
	return conf, err
}

func GetCofig() Configuration {
	return conf
}
