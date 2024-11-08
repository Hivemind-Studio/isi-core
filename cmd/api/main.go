package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	"github.com/Hivemind-Studio/isi-core/db"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"strings"
	"time"

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

type EmailData struct {
	Name    string
	Message string
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Inspirasi Satu",
		ErrorHandler: globalErrorHandler,
	})
	app.Use(cors.New())
	app.Use(compress.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.TimeoutMiddleware(5 * time.Second))
	app.Use(middleware.RequestIdMiddleware)

	config := configs.Init()

	api, _ := initApp(config)

	for _, r := range routerList(api) {
		r.RegisterRoutes(app)
	}

	log.Fatal(app.Listen(":3000"))
}

func dbInitConnection(cfg *configs.Config) *sqlx.DB {
	dbConf := cfg.Database
	dbConn := mysqlconn.Init(
		dbConf.Host,
		dbConf.Port,
		dbConf.Username,
		dbConf.Password,
		dbConf.DatabaseName)

	enableMigration := strings.EqualFold(dbConf.EnableDbMigration, "true")
	if enableMigration {
		db.InitMigration(dbConn.DB)
	}

	return dbConn
}

func globalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if customErr, ok := err.(*httperror.CustomError); ok {
		code = customErr.Code
	} else if fiberErr, ok := err.(*fiber.Error); ok {
		code = fiberErr.Code
	}

	log.Printf("Error: %v", err)

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": err.Error(),
		"code":    code,
	})
}
