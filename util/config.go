package util

import (
	"github.com/spf13/viper"
)

// Config store all configuration of the application
// the values read by viper from file or enviroment variables
type Config struct {
	Enviroment        string `mapstructure:"ENVIROMENT"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	HTTPAddressString string `mapstructure:"HTTP_ADDRESS_STRING"`
}

// LoadConfig read configuration from file conf.env or enviroment variables
func LoadConfig(configPath string) (config Config, err error) {
	viper.SetConfigName("conf")
	viper.SetConfigType("env")
	viper.AddConfigPath(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&config)
	return
}
