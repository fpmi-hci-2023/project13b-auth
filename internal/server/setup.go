package server

import (
	"fmt"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"github.com/fpmi-hci-2023/project13b-auth/config"
	"github.com/fpmi-hci-2023/project13b-auth/pkg/logger"
)

type Server struct {
	app *fiber.App
	log logger.Logger
}

func NewApp() Server {
	log := logger.NewLogger(config.DefaultWriter,
		config.LogInfo.Level,
		"auth-server")

	app := fiber.New(fiber.Config{
		Prefork:       config.GlobalConfig.FiberPrefork,
		ServerHeader:  "auth",
		CaseSensitive: false,
		ReadTimeout:   time.Second * 30,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {

			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			err = ctx.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("%d: %s", code, err))
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	})

	prometheus := fiberprometheus.New("auth")
	prometheus.RegisterAt(app, "/metrics")

	app.Use(
		//cors.New(cors.ConfigDefault),
		recover.New(),
		pprof.New(),
		prometheus.Middleware,
		logger.Middleware(
			logger.NewLogger(config.DefaultWriter,
				config.LogInfo.Level,
				"auth-httpserver"), nil,
		),
	)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("auth service healthy")
	})

	// Registering Swagger API
	app.Get("/swagger/*", swagger.HandlerDefault)

	return Server{
		app: app,
		log: log,
	}
}
