package router

import (
	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config, authController *controller.AuthController) {
	auth := api.Group("/auth")

	auth.Post("/register", authController.Register)
	auth.Post("/register/special-account", authController.RegisterSpecialAccount)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", authController.Logout)
	auth.Post("/verify-otp", authController.VerifyOTP)
	auth.Post("/complete-profile", authController.CompleteProfile)
	auth.Post("/resend-otp", authController.ResendOTP)
	auth.Post("/request-password-reset", authController.RequestPasswordReset)
	auth.Post("/reset-password", authController.ResetPassword)

	authProtected := auth.Group("")
	authProtected.Use(middleware.JWTAuth(cfg))
	authProtected.Post("/setup-pin", authController.SetupPIN)
}
