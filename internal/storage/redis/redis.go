package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
)

type Config struct {
	Host string
	Port string
}

func NewRedisClient(cfg *Config) (*redis.Client, error) {
	rc := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: "",
			DB:       0,
		},
	)

	// verify redis connection
	_, err := rc.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rc, nil
}
