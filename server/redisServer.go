package server

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func initRedis(config *viper.Viper) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.GetString("redis.addr"),
		Password: config.GetString("redis.password"), // no password set
		DB:       config.GetInt("redis.DB"),          // use default DB
	})
	return client
}
