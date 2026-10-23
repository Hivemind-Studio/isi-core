package configs

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
)

type Config struct {
	Database struct {
		Host              string `envconfig:"MYSQL_HOST"`
		Port              string `envconfig:"MYSQL_PORT"`
		Username          string `envconfig:"MYSQL_USER"`
		Password          string `envconfig:"MYSQL_PASSWORD"`
		DatabaseName      string `envconfig:"MYSQL_DB"`
		EnableDbMigration string `envconfig:"ENABLE_DB_MIGRATION"`
	}
	Mail struct {
		ApiKey string `envconfig:"MAIL_API_KEY"`
		Sender string `envconfig:"MAIL_SENDER"`
		Email  string `envconfig:"MAIL_EMAIL"`
		Url    string `envconfig:"MAIL_URL"`
	}
	GoogleConfig struct {
		ClientID       string `envconfig:"GOOGLE_CLIENT_ID"`
		ClientSecret   string `envconfig:"GOOGLE_CLIENT_SECRET"`
		RedirectURI    string `envconfig:"GOOGLE_REDIRECT_URI"`
		GoogleTokenUrl string `envconfig:"GOOGLE_TOKEN_URL"`
	}
	RedisConfig struct {
		Address   string `envconfig:"REDIS_ADDRESS"`
		Password  string `envconfig:"REDIS_PASSWORD"`
		DefaultDB string `envconfig:"REDIS_DEFAULT_DB"`
	}
	EnvConfig struct {
		Environment string `envconfig:"ENVIRONMENT"`
	}
}

func Init() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env not loaded, using default environment variables...", "error", err.Error())
	}

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		slog.Error("failed to load environment variable config.", "error", err.Error())
	}
	return &cfg
}
