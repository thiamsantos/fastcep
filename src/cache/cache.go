package cache

import (
	"os"

	"github.com/go-redis/redis"
)

// GetCredentials returns the redis credentials
func GetCredentials() *redis.Options {
	return &redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
}

// NewClient returns a client to the Redis Server
func NewClient(options *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(options)

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
