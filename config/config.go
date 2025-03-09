package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	Logger     LoggerConfig
	JWT        JWTConfig
	Cloudflare CloudflareConfig
}

type ServerConfig struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Driver   string
	URL      string
	MaxConns int
	MaxIdle  int
	Timeout  time.Duration
}

type LoggerConfig struct {
	Level string
	File  string
}

type JWTConfig struct {
	Secret         string
	ExpirationTime time.Duration
}

type CloudflareConfig struct {
	Token      string
	BucketName string
	AccountId  string
	AccessKey  string
	SecretKey  string
	Endpoint   string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	viper.SetDefault("server.address", ":8000")
	viper.SetDefault("server.readTimeout", 10*time.Second)
	viper.SetDefault("server.writeTimeout", 10*time.Second)
	viper.SetDefault("server.idleTimeout", 120*time.Second)

	viper.SetDefault("database.driver", "postgres")
	viper.SetDefault("database.url", "postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable")
	viper.SetDefault("database.maxConns", 20)
	viper.SetDefault("database.maxIdle", 5)
	viper.SetDefault("database.timeout", 5*time.Second)

	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.file", "")

	viper.SetDefault("jwt.secret", "your_secret_key_")
	viper.SetDefault("jwt.expirationTime", 25*time.Hour)

	viper.SetDefault("cloudflare.token", "your_cloudflare_token")
	viper.SetDefault("cloudflare.bucketName", "your_cloudflare_bn")
	viper.SetDefault("cloudflare.accountId", "your_cloudflare_accId")
	viper.SetDefault("cloudflare.accessKey", "your_cloudflare_ak")
	viper.SetDefault("cloudflare.secretKey", "your_cloudflare_sak")
	viper.SetDefault("cloudflare.endpoints", "endpoint")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil

}
