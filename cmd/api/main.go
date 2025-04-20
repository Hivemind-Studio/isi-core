package main

import (
	"github.com/Hivemind-Studio/isi-core/configs"
	"github.com/Hivemind-Studio/isi-core/db"
	"github.com/Hivemind-Studio/isi-core/internal/constant"
	"github.com/Hivemind-Studio/isi-core/pkg/googleoauth2"
	"github.com/Hivemind-Studio/isi-core/pkg/httperror"
	"github.com/Hivemind-Studio/isi-core/pkg/mail"
	"github.com/Hivemind-Studio/isi-core/pkg/middleware"
	"github.com/Hivemind-Studio/isi-core/pkg/mysqlconn"
	redisutils "github.com/Hivemind-Studio/isi-core/pkg/redis"
	"github.com/Hivemind-Studio/isi-core/pkg/session"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"golang.org/x/oauth2"

	"os"
	"strconv"
	"strings"
	"time"
)

var requestCount *prometheus.CounterVec

func main() {
	initLogger()

	app := fiber.New(fiber.Config{
		AppName:      "Inspirasi Satu",
		ErrorHandler: globalErrorHandler,
	})

	initMetrics()

	// Middleware to count requests with method, route, and status
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()
		status := c.Response().StatusCode()

		route := c.Route().Path
		if route == "" {
			route = c.Path()
		}

		requestCount.WithLabelValues(c.Method(), route, strconv.Itoa(status)).Inc()
		log.Info().
			Str("method", c.Method()).
			Str("route", route).
			Int("status", status).
			Dur("duration", time.Since(start)).
			Msg("Request handled")

		return err
	})

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Use(compress.New())
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

		err := c.Next()

		status := c.Response().StatusCode()

		// If the status code is not 2xx, log as an error
		if status < 200 || status >= 300 {
			log.Error().
				Str("request_id", c.Locals("request_id").(string)).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Str("origin", origin).
				Str("cookie", cookie).
				Str("user_agent", userAgent).
				Str("ip", c.IP()).
				Int("status", status).
				Str("request_id", c.Locals("request_id").(string)).
				Bytes("body", c.Response().Body())
		} else {
			log.Info().
				Str("request_id", c.Locals("request_id").(string)).
				Str("method", c.Method()).
				Str("path", c.Path()).
				Str("origin", origin).
				Str("cookie", cookie).
				Str("user_agent", userAgent).
				Str("ip", c.IP()).
				Int("status", status).
				Bytes("body", c.Response().Body())
		}

		return err
	})

	config := configs.Init()
	sessionManager := initSessionManager(config)
	api, _ := initApp(config, sessionManager)

	for _, r := range routerList(api) {
		r.RegisterRoutes(app, sessionManager)
	}

	log.Fatal().Err(app.Listen(os.Getenv("APP_PORT"))).Msg("Fiber server exited")
}

func initMetrics() {
	requestCountVec := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests by method, route, and status",
		},
		[]string{"method", "route", "status"},
	)

	if err := prometheus.Register(requestCountVec); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			requestCount = are.ExistingCollector.(*prometheus.CounterVec)
		} else {
			panic(err)
		}
	} else {
		requestCount = requestCountVec
	}
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
		APIKey: mailConfig.ApiKey,
		Sender: mailConfig.Sender,
		Email:  mailConfig.Email,
		Url:    mailConfig.Url,
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

	log.Error().Err(err).Int("status_code", code).Msg("Unhandled error occurred")

	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": err.Error(),
		"code":    code,
	})
}

func initLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}
