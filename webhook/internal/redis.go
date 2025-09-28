package internal

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		client: client,
	}
}

func (r *RedisService) StoreSuscriptor(token string, suscriptor map[string]string) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(suscriptor)
	if err != nil {
		return err
	}

	key := "suscriptor:" + token
	return r.client.Set(ctx, key, jsonData, 0).Err()
}

func (r *RedisService) GetSuscriptor(token string) (map[string]string, error) {
	ctx := context.Background()

	key := "suscriptor:" + token
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var suscriptor map[string]string
	err = json.Unmarshal([]byte(val), &suscriptor)
	if err != nil {
		return nil, err
	}

	return suscriptor, nil
}
