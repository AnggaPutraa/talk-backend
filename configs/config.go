package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             int    `mapstructure:"DB_PORT"`
	DBUsername         string `mapstructure:"DB_USERNAME"`
	DBPassword         string `mapstructure:"DB_PASSWORD"`
	DBDatabase         string `mapstructure:"DB_DATABASE"`
	AccessTokenSecret  string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	ServerAddress      string `mapstructure:"SERVER_ADDRESS"`
	DBUrl              string
}

func LoadConfig() (config Config) {
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err = viper.Unmarshal(&config); err != nil {
		log.Fatal("Can't load environment variable: ", err)
	}
	config.DBUrl = fmt.Sprintf("postgres://%s:%s@%s:%d?sslmode=disable",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
	)
	return config
}
