package redis

import (
	"context"
	"encoding/json"
	"time"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(addr, password string, db int) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Redis{client: rdb}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, timeout time.Duration) error {
	j, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	if _, err = r.client.Set(ctx, key, j, timeout).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	b, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if value != nil {
		if err := json.Unmarshal(b, value); err != nil {
			return err
		}
	}
	return nil
}

func (r *Redis) GetByKeyPattern(ctx context.Context, key string) ([]string, error) {
	return r.client.Keys(ctx, key).Result()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) DelPattern(ctx context.Context, key string) error {
	keys, err := r.client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}
	return r.client.Del(ctx, keys...).Err()
}

func (r *Redis) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}

func (r *Redis) IsKeyNotFound(err error) bool {
	return err == redis.Nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
