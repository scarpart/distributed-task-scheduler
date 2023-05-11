package util 

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_DRIVER string `mapstructure:"DB_DRIVER"`
	DB_SOURCE string `mapstructure:"DB_SOURCE"`
	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
	LB_CONN_TIMEOUT int32 `mapstructure:"LB_CONN_TIMEOUT"`
	LB_CLIENT_MAX_CONNS int32 `mapstructure:"LB_MAX_CONNS"`
	LB_HEALTH_CHECK_INTERVAL int32 `mapstructure:"LB_HEALTH_CHECK_INTERVAL"`
}

// Read constants from a config file using viper   
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return 
	}
	
	err = viper.Unmarshal(&config)
	return
}
