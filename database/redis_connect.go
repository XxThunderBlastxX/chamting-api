package database

import (
	"github.com/go-redis/redis/v9"
	"os"
)

// RedisConnect is a method which connects and returns redis client
func RedisConnect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASS"), DB: 0})

	return rdb
}
