package redis

import (
	"context"
	"encoding/json"
	"go-service-demo/internal/model"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDatabase struct {
	redis *redis.Client
}

func NewRedisClient() *RedisDatabase {
	return &RedisDatabase{}
}

func (r *RedisDatabase) Connect() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	r.redis = rdb
	return r.Ping()
}

func (r *RedisDatabase) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ping := r.redis.Ping(ctx)
	if ping.Err() != nil {
		return ping.Err()
	}
	return nil
}

func (r *RedisDatabase) Close() error {
	return r.redis.Close()
}

func (r *RedisDatabase) Set(key string, value string, expiredTimeInMs int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.Set(ctx, key, value, time.Duration(expiredTimeInMs)*time.Millisecond).Err()
}

func (r *RedisDatabase) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return r.redis.Get(ctx, key).Result()
}

func (r *RedisDatabase) SaveUserToRedis(key string, user model.User) {
	bytes, _ := json.Marshal(user)
	err := r.Set(key, string(bytes), 300000)
	if err != nil {
		log.Println("Error when set redis: " + err.Error())
	}
}
