package component

import (
	"entry_task/internal/config"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	Redis *redis.Client
}

func InitRedis() *Cache {
	var cache *Cache
	conf := config.Get()

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       0,
	})

	cache = &Cache{
		Redis: client,
	}

	fmt.Println("Redis Client Initialized")

	return cache
}
