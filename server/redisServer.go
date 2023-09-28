package server

import (
	"github.com/redis/go-redis/v9"
)

// func initRedis(config *viper.Viper) *redis.Client {
func initRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
