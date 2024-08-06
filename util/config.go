package util

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"

)

type Config struct {
	DBUser                 string `mapstructure:"DB_USER"`
	DBPassword             string `mapstructure:"DB_PASSWORD"`
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBName                 string `mapstructure:"DB_NAME"`
	SecretKey              string `mapstructure:"SECRET_KEY"`
	LogtoEndpoint          string `mapstructure:"LOGTO_ENDPOINT"`
	LogtoApplicationId     string `mapstructure:"LOGTO_APPLICATION_ID"`
	LogtoApplicationSecret string `mapstructure:"LOGTO_APPLICATION_SECRET"`
	CloudinaryApiKey       string `mapstructure:"CLOUDINARY_API_KEY"`
	CloudinaryApiSecret    string `mapstructure:"CLOUDINARY_API_SECRET"`
	CloudinaryCloudName    string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	RedisUrl               string `mapstructure:"REDIS_URL"`
}

var (
	config Config
	once   sync.Once
	loaded bool
)

func loadConfig() (Config, error) {
	var err error

	once.Do(func() {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		viper.AutomaticEnv()

		if err = viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				fmt.Println("Config file not found, relying on environment variables")
			} else {
				return
			}
		}

		viper.BindEnv("DB_USER")
		viper.BindEnv("DB_PASSWORD")
		viper.BindEnv("DB_HOST")
		viper.BindEnv("DB_PORT")
		viper.BindEnv("DB_NAME")
		viper.BindEnv("SECRET_KEY")
		viper.BindEnv("LOGTO_ENDPOINT")
		viper.BindEnv("LOGTO_APPLICATION_ID")
		viper.BindEnv("LOGTO_APPLICATION_SECRET")
		viper.BindEnv("CLOUDINARY_API_KEY")
		viper.BindEnv("CLOUDINARY_API_SECRET")
		viper.BindEnv("CLOUDINARY_CLOUD_NAME")
		viper.BindEnv("REDIS_URL")

		err = viper.Unmarshal(&config)
		if err == nil {
			loaded = true
		}
	})

	return config, err
}

func GetConfig() (Config, error) {
	if !loaded {
		return loadConfig()
	}
	return config, nil
}
