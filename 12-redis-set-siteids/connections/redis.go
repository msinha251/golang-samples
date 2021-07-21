package connections

import (
	"github.com/go-redis/redis"
)

func RedisConnect() (*redis.Client, error) {
	rds := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rds, nil
}
