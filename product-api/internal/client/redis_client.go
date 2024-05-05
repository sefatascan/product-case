package client

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"product-api/internal/config"
)

type RedisClient interface {
	Subscribe(callback func(payload string))
}

type RedisClientImpl struct {
	applicationConfig config.ApplicationConfigManager
	client            *redis.Client
}

func (r *RedisClientImpl) Subscribe(callback func(payload string)) {
	subscribe := r.client.Subscribe(context.Background(), r.applicationConfig.Redis.EventName)

	_, err := subscribe.Receive(context.Background())

	if err != nil {
		log.Error(err)
		panic(err)
	}

	channel := subscribe.Channel()

	go func() {
		for msg := range channel {
			callback(msg.Payload)
		}
	}()
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
