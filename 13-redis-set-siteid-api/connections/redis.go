package connections

import (
	"os"

	"github.com/go-redis/redis"
)

func RedisConnect() (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rds, nil
}
