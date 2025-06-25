package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

func InitRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return rdb
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client: client}
}

func (r *RedisStore) Save(shortUrl, longUrl, userId, customAlias string) error {
	ctx := context.Background()
	return r.client.Set(ctx, shortUrl, longUrl, time.Hour*24).Err()
}

func (r *RedisStore) Get(shortUrl string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, shortUrl).Result()
}

func (r *RedisStore) Exists(shortUrl string) bool {
	_, err := r.Get(shortUrl)
	return err == nil
}
