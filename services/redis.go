package services

import "github.com/go-redis/redis"

var RedisClient *redis.Client

func InitRedis(address, password string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: 		address,
		Password: 	"",
		DB: 		0,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}
