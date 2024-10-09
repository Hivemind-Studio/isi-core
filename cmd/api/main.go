package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
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
	app.Use(logger.New())
	app.Use(recover.New())

	config := configs.Init()

	api, _ := initApp(config)

	app.Group("/api/v1", middleware.AuthMiddleware())

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
