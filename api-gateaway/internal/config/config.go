package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	Environment          string   `mapstructure:"ENVIRONMENT"`
	HttpServerAddress    string   `mapstructure:"HTTP_SERVER_ADDRESS"`
	AuthServerAddress    string   `mapstructure:"AUTH_SERVER_ADDRESS"`
	UserServerAddress    string   `mapstructure:"USER_SERVER_ADDRESS"`
	CarsServerAddress    string   `mapstructure:"CARS_SERVER_ADDRESS"`
	OrderServerAddress   string   `mapstructure:"ORDER_SERVER_ADDRESS"`
	PaymentServerAddress string   `mapstructure:"PAYMENT_SERVER_ADDRESS"`
	AllowedOrigins       []string `mapstructure:"ALLOWED_ORIGINS"`
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
