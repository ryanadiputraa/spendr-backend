package config

import "github.com/spf13/viper"

type Config struct {
	*Server
	*Postgres
}

type Server struct {
	Port string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func LoadConfig(configType, filePath string) (*Config, error) {
	viper.SetConfigType(configType)
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config *Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return config, nil
}
