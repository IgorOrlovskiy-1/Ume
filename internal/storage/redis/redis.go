package redis

import (
	"github.com/go-redis/redis/v8"
    "context"
    "fmt"
)

type RedisClient struct {
	redisClient *redis.Client
}

func NewRedisClient(ctx context.Context) (*RedisClient, error) {
	op := "storage.redis.NewRedisClient"

    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   0, 
    })

    _, err := client.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &RedisClient{redisClient: client}, nil
}