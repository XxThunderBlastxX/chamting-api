package database

import (
	"github.com/go-redis/redis/v8"
	"os"
)

// RedisConnect is a method which connects and returns redis client
func RedisConnect(db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASS"), DB: db})

	return rdb
}
