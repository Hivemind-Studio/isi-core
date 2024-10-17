package main

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/configs"
	"github.com/Hivemind-Studio/isi-core/pkg/redis"
	//"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/mysqlconn"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{AppName: "Awesome Project"})
	app.Use(cors.New())
	app.Use(compress.New())
	//app.Use(utils.LoggerMiddleware(utils.NewLogger()))
	app.Use(logger.New())
	app.Use(recover.New())

	config := configs.Init()

	err := redis.InitRedis("localhost:6379", "", 0)

	if err != nil {
		logger.New(logger.Config{
			Format: "Failed to initialize Redis",
		})
	}

	client := redis.GetRedisClient()

	err = client.Set(context.Background(), "key", "value", 0).Err()
	if err != nil {
		log.Fatalf("Failed to set key: %v", err)
	}

	api, _ := initApp(config)

	for _, r := range routerList(api) {
		r.RegisterRoutes(app)
	}

	log.Fatal(app.Listen(":3000"))
}

func dbInitConnection(cfg *configs.Config) *sqlx.DB {
	dbConf := cfg.Database
	return mysqlconn.Init(
		dbConf.Host,
		dbConf.Port,
		dbConf.Username,
		dbConf.Password,
		dbConf.DatabaseName)
}
