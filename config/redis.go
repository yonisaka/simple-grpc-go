package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func ConnectRedis() {
	rdHost := viper.GetString(`redis.host`)
	rdPort := viper.GetString(`redis.port`)
	client := redis.NewClient(&redis.Options{
		Addr: rdHost+rdPort,
	})
	
	RedisClient = client
}