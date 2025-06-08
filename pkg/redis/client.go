package redis

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
)

type RedisConfig struct {
	Host     string
	Port     int64
	Password string
	Database int64
}

func GetRedisConnection(config RedisConfig) (*redis.Client, error) {
	if Client != nil {
		return Client, nil
	}

	ctx := context.Background()

	Client = redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + strconv.Itoa(int(config.Port)),
		Password: config.Password,
		DB:       int(config.Database),
	})

	if err := Client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return Client, nil
}
