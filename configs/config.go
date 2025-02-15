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
		Host      string `envconfig:"MAIL_SMTP_HOST"`
		Port      string `envconfig:"MAIL_SMTP_PORT"`
		Username  string `envconfig:"MAIL_SMTP_USERNAME"`
		Password  string `envconfig:"MAIL_SMTP_PASSWORD"`
		NameFrom  string `envconfig:"MAIL_SENDER_NAME"`
		EmailFrom string `envconfig:"MAIL_AUTH_EMAIL"`
	}
	GoogleConfig struct {
		ClientID     string `envconfig:"GOOGLE_CLIENT_ID"`
		ClientSecret string `envconfig:"GOOGLE_CLIENT_SECRET"`
		RedirectURI  string `envconfig:"GOOGLE_REDIRECT_URI"`
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
