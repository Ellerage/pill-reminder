package configs

import (
	"log"
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	BOT_TOKEN     string `env:"BOT_TOKEN" env-default:""`
	MONGO_URL     string `env:"MONGO_URL" env-default:""`
	MONGO_DB_NAME string `env:"MONGO_DB_NAME" env-default:""`
	TIMEZONE      string `env:"TIMEZONE" env-default:""`
}

var (
	cfg  *Config
	once sync.Once
)

func InitConfig() *Config {
	once.Do(func() {
		err := godotenv.Load(".env")

		if err != nil {
			slog.Error("No .env file found (optional)")
		}

		cfg = &Config{}
		if err := cleanenv.ReadEnv(cfg); err != nil {
			slog.Error("Failed to read env", slog.Any("error", err))
		}

		log.Println("Init config")
	})

	return cfg
}

func GetConfig() *Config {
	if cfg == nil {
		log.Fatal("Config not initialized. Call config.InitConfig() first.")
	}
	return cfg
}
