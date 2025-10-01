package internal

import (
	"context"
	"encoding/json"
	"time"

	"cards/models"

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

func (r *RedisService) StoreUser(ctx context.Context, token string, user models.User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := "user:" + token
	return r.client.Set(ctx, key, userJSON, 24*time.Hour).Err()
}

func (r *RedisService) GetUser(ctx context.Context, token string) (*models.User, error) {
	key := "user:" + token
	userJSON, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RedisService) StoreRequest(ctx context.Context, uuid string, data models.RequestData) error {
	requestJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := "request:" + uuid
	return r.client.Set(ctx, key, requestJSON, 24*time.Hour).Err()
}

func (r *RedisService) GetRequest(ctx context.Context, uuid string) (*models.RequestData, error) {
	key := "request:" + uuid
	requestJSON, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var data models.RequestData
	err = json.Unmarshal([]byte(requestJSON), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *RedisService) DeleteRequest(ctx context.Context, uuid string) error {
	key := "request:" + uuid
	return r.client.Del(ctx, key).Err()
}

func (r *RedisService) GetAllUserKeys(ctx context.Context) ([]string, error) {
	return r.client.Keys(ctx, "user:*").Result()
}
