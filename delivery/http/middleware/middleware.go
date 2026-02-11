package middleware

import (
	"runtime/debug"
	"time"

	"iam-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

type Middleware struct {
	config *config.Config
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *Middleware {
	return &Middleware{
		config: cfg,
		logger: logger,
	}
}

func (m *Middleware) Setup(app *fiber.App) {
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			m.logger.Error("panic recovered",
				zap.Any("error", e),
				zap.String("stack", string(debug.Stack())),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.String("request_id", c.GetRespHeader("X-Request-ID")),
			)
		},
	}))

	app.Use(requestid.New())

	app.Use(RequestContext())

	app.Use(RequestLogger(m.logger))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	if m.config.IsProduction() {
		app.Use(limiter.New(limiter.Config{
			Max:               10,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.SlidingWindow{},
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error":   "too many requests",
				})
			},
		}))
	}
}
