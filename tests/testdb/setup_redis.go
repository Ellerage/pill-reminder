package testsdb

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var RedisPort int
var RedisClient *redis.Client

func SetupRedis() func() {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7.0",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}
	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(err)
	}

	host, err := redisC.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port.Port()),
	})

	err = client.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}

	RedisPort = port.Int()
	RedisClient = client

	teardown := func() {
		_ = client.Close()
		_ = redisC.Terminate(ctx)
	}

	return teardown
}
