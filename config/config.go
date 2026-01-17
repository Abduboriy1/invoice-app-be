// config/config.go
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	Jira     JiraConfig
	Square   SquareConfig
	Redis    RedisConfig
}

type ServerConfig struct {
	Port         int
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

type AuthConfig struct {
	JWTSecret     string
	TokenDuration time.Duration
}

type JiraConfig struct {
	BaseURL string
	Enabled bool
}

type SquareConfig struct {
	AccessToken string
	Environment string // sandbox or production
	Enabled     bool
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.maxconns", 25)

	// Read config file (if exists)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	// Bind environment variables - MUST BE AFTER ReadInConfig
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Explicitly bind important vars
	viper.BindEnv("database.host", "APP_DATABASE_HOST")
	viper.BindEnv("database.port", "APP_DATABASE_PORT")
	viper.BindEnv("database.user", "APP_DATABASE_USER")
	viper.BindEnv("database.password", "APP_DATABASE_PASSWORD")
	viper.BindEnv("database.dbname", "APP_DATABASE_DBNAME")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}
