package configs

import (
	"log"
	"log/slog"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	BOT_TOKEN                  string `env:"BOT_TOKEN" env-default:""`
	MONGO_URL                  string `env:"MONGO_URL" env-default:"mongodb://localhost:27017/pill-reminder"`
	MONGO_DB_NAME              string `env:"MONGO_DB_NAME" env-default:"pill-reminder"`
	TIMEZONE                   string `env:"TIMEZONE" env-default:"UTC"`
	REDIS_URL                  string `env:"REDIS_URL" env-default:"127.0.0.1"`
	REDIS_PASSWORD             string `env:"REDIS_PASSWORD" env-default:""`
	ASYNCQ_DB                  int    `env:"ASYNCQ_DB" env-default:"0"`
	REDIS_PORT                 int    `env:"REDIS_PORT" env-default:"6379"`
	REMINDER_DB                int    `env:"REMINDER_DB" env-default:"1"`
	MONGO_INITDB_ROOT_USERNAME string `env:"MONGO_INITDB_ROOT_USERNAME" env-default:"root"`
	MONGO_INITDB_ROOT_PASSWORD string `env:"MONGO_INITDB_ROOT_PASSWORD" env-default:"admin"`
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
		} else {
			slog.Info("Loaded .env file")
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
