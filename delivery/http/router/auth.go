package router

import (
	"time"

	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config, authController *controller.AuthController) {
	auth := api.Group("/auth")

	if !cfg.IsDevelopment() {
		auth.Use(limiter.New(limiter.Config{
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

	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/verify-otp", authController.VerifyOTP)
	auth.Post("/complete-profile", authController.CompleteProfile)
	auth.Post("/resend-otp", authController.ResendOTP)
	auth.Post("/request-password-reset", authController.RequestPasswordReset)
	auth.Post("/reset-password", authController.ResetPassword)

	authProtected := auth.Group("")
	authProtected.Use(middleware.JWTAuth(cfg))
	authProtected.Post("/logout", authController.Logout)
	authProtected.Post("/setup-pin", authController.SetupPIN)

	authAdmin := auth.Group("")
	authAdmin.Use(middleware.JWTAuth(cfg))
	authAdmin.Use(middleware.RequirePlatformAdmin())
	authAdmin.Post("/register/special-account", authController.RegisterSpecialAccount)

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
	registrations.Post("/:id/complete", authController.CompleteRegistration)
}
