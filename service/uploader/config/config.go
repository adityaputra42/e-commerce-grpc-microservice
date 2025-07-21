package config

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	SupabaseUrl string `mapstructure:"SUPABASE_URL"`
	SupabaseKey string `mapstructure:"SUPABASE_ANON_PUBLIC_KEY"`
	Bucket      string `mapstructure:"BUCKET"`
}

var conf Configuration

func LoadConfig(path string) (config Configuration, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigFile("upload.env")

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
