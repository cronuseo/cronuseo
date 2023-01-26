package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/shashimalcse/cronuseo/internal/entity"
)

type redisCache struct {
	host     string
	db       int
	expire   time.Duration
	password string
}

func NewRedisCache(host string, db int, expire time.Duration, password string) PermissionCache {

	return &redisCache{
		host:     host,
		db:       db,
		expire:   expire,
		password: password,
	}
}

func (c *redisCache) getClient() *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     c.host,
		Password: c.password,
		DB:       c.db,
	})
}

func (c *redisCache) Get(context context.Context, key entity.Tuple) (string, error) {

	client := c.getClient()

	k := key.String()
	value, _ := client.Get(context, k).Result()

	return value, nil
}

func (c *redisCache) Set(context context.Context, key entity.Tuple, value string) error {

	client := c.getClient()

	k := key.String()
	client.Set(context, k, value, c.expire*time.Second)

	return nil
}

func (c *redisCache) FlushAll(context context.Context) error {

	client := c.getClient()
	client.FlushAll(context)
	return nil
}

func (c *redisCache) SetAPIKey(context context.Context, key string, value string) error {

	client := c.getClient()

	client.Set(context, key, value, c.expire*time.Second)

	return nil

}

func (c *redisCache) GetAPIKey(context context.Context, key string) (string, error) {

	client := c.getClient()

	value, _ := client.Get(context, key).Result()

	return value, nil
}

func (c *redisCache) DeleteAPIKey(context context.Context) error {

	client := c.getClient()
	client.Del(context, "API_KEY")
	return nil
}
