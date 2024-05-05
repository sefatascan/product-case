package client

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"product-config-api/internal/config"
	"time"
)

type RedisClient interface {
	PublishEvent() error
}

type RedisClientImpl struct {
	applicationConfig config.ApplicationConfigManager
	client            *redis.Client
}

func (r *RedisClientImpl) PublishEvent() error {
	err := r.client.Publish(context.Background(), r.applicationConfig.Redis.EventName, time.Now()).Err()
	return err

}

func NewRedisClient(applicationConfig config.ApplicationConfigManager) RedisClient {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", applicationConfig.Redis.Host, applicationConfig.Redis.Port),
		Password: "",
		DB:       0,
	})
	return &RedisClientImpl{
		applicationConfig: applicationConfig,
		client:            redisClient,
	}
}
