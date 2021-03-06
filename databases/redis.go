package databases

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

type RedistInterface interface {
	RedisConnection() (*redis.Client, error)
}

type RedisHandler struct {
	h RedistInterface
}

type Redis struct {
	*redis.Client
}

func NewRedisHandler() (*RedisHandler, error) {
	return new(RedisHandler), nil
}

func (h *RedisHandler) RedisConnection() (*Redis, error) {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	var client Redis
	client.Client = redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()

	return &client, err
}
