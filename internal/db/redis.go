package db

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type ConnectRedisOptions struct {
	Addr     string
	Port     int
	DB       int
	Password string
}

func ConnectRedis(ctx context.Context, options ConnectRedisOptions) *redis.Client {
	var rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", options.Addr, options.Port),
		Password: options.Password,
		DB:       options.DB,
	})

	status := rdb.Ping(ctx)

	slog.Info("Redis connect: " + status.String())

	return rdb
}
