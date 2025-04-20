package redis

import (
	"context"
	logger2 "github.com/Hivemind-Studio/isi-core/pkg/logger"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient  *redis.Client
	redisOptions *redis.Options
)

func InitRedis(addr string, password string, db int) *redis.Client {
	logger := logger2.GetLogger()

	// Set up Redis options
	redisOptions = &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	}

	// Create a new Redis client
	redisClient = redis.NewClient(redisOptions)

	// Test the connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error().Err(err).Msg("redis ping failed")
		return nil
	}

	logger.Info().Msg("redis ping success")
	return redisClient
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatal("Redis client is not initialized")
	}
	return redisClient
}
