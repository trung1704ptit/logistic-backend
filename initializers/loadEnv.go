package initializers

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST_DB"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	UploadFilePath string `mapstructure:"UPLOAD_FILE_PATH"`
}

func LoadConfig(path string) (config Config, err error) {
	// Default environment is dev, you can change this based on your needs
	// env APP_ENV=dev go run main.go for dev command
	env := os.Getenv("APP_ENV")

	// Determine which config file to load based on the environment
	var configFile string
	if env == "dev" {
		configFile = "app.dev.env"
	} else {
		configFile = "app.prod.env"
	}

	// Set the path and file for viper
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(configFile)

	// Automatically read environment variables
	viper.AutomaticEnv()

	// Read in the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Unmarshal the config into the struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	return config, nil
}
