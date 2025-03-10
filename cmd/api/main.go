package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	"github.com/Hivemind-Studio/isi-core/db"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/googleoauth2"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/mail"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	redisutils "github.com/Hivemind-Studio/isi-core/pkg/redis"
	"github.com/Hivemind-Studio/isi-core/pkg/session"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"golang.org/x/oauth2"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Hivemind-Studio/isi-core/pkg/mysqlconn"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"log"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Inspirasi Satu",
		ErrorHandler: globalErrorHandler,
	})
	app.Use(compress.New())
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(middleware.TimeoutMiddleware(5 * time.Second))
	app.Use(middleware.RequestIdMiddleware)
	app.Use(middleware.BodyLimit("3MB"))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://dashboard.inspirasisatu.com, https://backoffice.inspirasisatu.com",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,OPTIONS,PUT,PATCH,DELETE",
		AllowHeaders:     "Content-Type, Authorization," + constant.APP_ORIGIN_HEADER,
	}))

	app.Use(func(c *fiber.Ctx) error {
		origin := c.Get("Origin")
		cookie := string(c.Request().Header.Peek("Cookie"))
		userAgent := string(c.Request().Header.Peek("User-Agent"))

		log.Printf("Incoming Request: method=%s path=%s origin=%s cookie=(%s) userAgent=%s ip=%s",
			c.Method(), c.Path(), origin, cookie, userAgent, c.IP(),
		)

		return c.Next()
	})

	config := configs.Init()

	sessionManager := initSessionManager(config)
	api, _ := initApp(config, sessionManager)

	for _, r := range routerList(api) {
		r.RegisterRoutes(app, sessionManager)
	}

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
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

func initEmailClient(cfg *configs.Config) *mail.EmailClient {
	mailConfig := cfg.Mail
	emailClient := mail.NewEmailClient(&mail.EmailConfig{
		Host:        mailConfig.Host,
		Port:        mailConfig.Port,
		Username:    mailConfig.Username,
		Password:    mailConfig.Password,
		SenderEmail: mailConfig.EmailFrom,
		SenderName:  mailConfig.NameFrom,
	})

	return emailClient
}

func initGoogleOauthClient(cfg *configs.Config) *oauth2.Config {
	googleOauthConfig := cfg.GoogleConfig
	return googleoauth2.NewGoogleOauth(&googleoauth2.OauthConfig{
		ClientID:     googleOauthConfig.ClientID,
		ClientSecret: googleOauthConfig.ClientSecret,
		RedirectURL:  googleOauthConfig.RedirectURI,
	})
}

func initSessionManager(cfg *configs.Config) *session.SessionManager {
	num, _ := strconv.ParseInt(cfg.RedisConfig.DefaultDB, 10, 32)
	redisClient := redisutils.InitRedis(cfg.RedisConfig.Address, cfg.RedisConfig.Password, int(num))
	return session.NewSessionManager(redisClient)
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
