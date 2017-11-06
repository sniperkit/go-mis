package services

import (
	"fmt"
	"log"
	"strconv"

	"bitbucket.org/go-mis/config"
	"github.com/go-redis/redis"
)

// RedisClient - store redis client driver and additional information
type RedisClient struct {
	*redis.Client
}

// NewClientRedis - instance new client redis
func NewClientRedis() (*RedisClient, error) {
	redisClient := new(RedisClient)
	port := strconv.Itoa(config.Configuration.Redis.Port)
	fmt.Println("Redis port: ", port)
	client := redis.NewClient(&redis.Options{
		Addr:     config.Configuration.Redis.Address + ":" + port,
		Password: config.Configuration.Redis.Password,
		DB:       config.Configuration.Redis.Db,
	})
	if _, err := client.Ping().Result(); err != nil {
		log.Println("[ERROR] Unable to connect to redis", err)
		panic(err)
	}
	redisClient.Client = client
	return redisClient, nil
}
