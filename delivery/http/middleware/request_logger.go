package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger(zapLogger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := c.GetRespHeader("X-Request-ID")

		zapLogger.Info("request_started",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.String("query", string(c.Request().URI().QueryString())),
		)

		err := c.Next()

		zapLogger.Info("request_completed",
			zap.String("request_id", requestID),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.Duration("duration", time.Since(start)),
			zap.Bool("success", c.Response().StatusCode() < 400),
		)

		return err
	}
}
