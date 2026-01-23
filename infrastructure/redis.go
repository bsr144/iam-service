package infrastructure

import (
	"context"
	"fmt"
	"log"

	"iam-service/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetAddress(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	log.Printf("Connected to Redis at %s\n", cfg.GetAddress())

	return client, nil
}
