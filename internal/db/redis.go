package db

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type ConnectRedisOptions struct {
	Addr     string
	Port     int
	DB       int
	Password string
}

func ConnectRedis(ctx context.Context, options ConnectRedisOptions) *redis.Client {
	addr := net.JoinHostPort(options.Addr, strconv.Itoa(options.Port))

	var rdb = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     options.Password,
		DB:           options.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		PoolSize:     5,
		MinIdleConns: 2,
	})

	status := rdb.Ping(ctx)

	slog.Info(fmt.Sprintf("Redis connect: %s. With: Addr: %s. DB: %v", status.String(), addr, options.DB))

	return rdb
}
