package router

import (
	"time"

	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"
	"iam-service/iam/auth/contract"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config, authController *controller.AuthController, blacklistStore contract.TokenBlacklistStore) {
	auth := api.Group("/auth")
	auth.Use(middleware.JWTAuth(cfg, blacklistStore))
	auth.Post("/logout", authController.Logout)
	auth.Post("/logout-all", authController.LogoutAll)

	registrations := api.Group("/registrations")
	if !cfg.IsDevelopment() {
		registrations.Use(limiter.New(limiter.Config{
			Max:               10,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.SlidingWindow{},
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error":   "too many requests, please try again later",
				})
			},
		}))
	}
	registrations.Post("", authController.InitiateRegistration)
	registrations.Post("/:id/verify-otp", authController.VerifyRegistrationOTP)
	registrations.Post("/:id/resend-otp", authController.ResendRegistrationOTP)
	registrations.Get("/:id/status", authController.GetRegistrationStatus)
	registrations.Post("/:id/set-password", authController.SetPassword)
	registrations.Post("/:id/complete-profile", authController.CompleteProfileRegistration)
	registrations.Post("/:id/complete", authController.CompleteRegistration)

	login := api.Group("/login")
	if !cfg.IsDevelopment() {
		login.Use(limiter.New(limiter.Config{
			Max:               10,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.SlidingWindow{},
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error":   "too many requests, please try again later",
				})
			},
		}))
	}
	login.Post("", authController.InitiateLogin)
	login.Post("/:id/verify-otp", authController.VerifyLoginOTP)
	login.Post("/:id/resend-otp", authController.ResendLoginOTP)
	login.Get("/:id/status", authController.GetLoginStatus)
}
