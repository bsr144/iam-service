package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"
	"iam-service/iam/auth"
	"iam-service/iam/health"
	"iam-service/iam/role"
	"iam-service/iam/user"
	"iam-service/impl/mailer"
	"iam-service/impl/postgres"
	"iam-service/infrastructure"
	"iam-service/pkg/logger"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	app    *fiber.App
	config *config.Config
	logger *zap.Logger

	healthUsecase health.Usecase
	authUsecase   auth.Usecase
	roleUsecase   role.Usecase
	userUsecase   user.Usecase
}

func NewServer(cfg *config.Config) *Server {
	zapLogger, _ := logger.NewZapLogger(cfg.App.Environment)

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		AppName:      cfg.App.Name,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: defaultErrorHandler,
	})

	postgresDB, err := infrastructure.NewPostgres(cfg.Infra.Postgres, zapLogger)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}

	authUserRepo := postgres.NewUserRepository(postgresDB)
	userProfileRepo := postgres.NewUserProfileRepository(postgresDB)
	userCredentialsRepo := postgres.NewUserCredentialsRepository(postgresDB)
	userSecurityRepo := postgres.NewUserSecurityRepository(postgresDB)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(postgresDB)
	tenantRepo := postgres.NewTenantRepository(postgresDB)
	userActivationTrackingRepo := postgres.NewUserActivationTrackingRepository(postgresDB)
	roleRepo := postgres.NewRoleRepository(postgresDB)

	emailService := mailer.NewEmailService()

	healthUsecase := health.NewUsecase()
	authUsecase := auth.NewUsecase(
		postgresDB,
		cfg,
		authUserRepo,
		userProfileRepo,
		userCredentialsRepo,
		userSecurityRepo,
		emailVerificationRepo,
		tenantRepo,
		userActivationTrackingRepo,
		roleRepo,
		emailService,
	)
	roleUsecase := role.NewUsecase(
		postgresDB,
		cfg,
		tenantRepo,
		roleRepo,
	)
	userUsecase := user.NewUsecase(
		postgresDB,
		cfg,
		authUserRepo,
		userProfileRepo,
		userCredentialsRepo,
		userSecurityRepo,
		tenantRepo,
		roleRepo,
		userActivationTrackingRepo,
	)

	server := &Server{
		app:           app,
		config:        cfg,
		logger:        zapLogger,
		healthUsecase: healthUsecase,
		authUsecase:   authUsecase,
		roleUsecase:   roleUsecase,
		userUsecase:   userUsecase,
	}

	server.setupMiddleware()
	server.setupRoutes()

	return server
}

func (s *Server) setupMiddleware() {
	mw := middleware.New(s.config, s.logger)
	mw.Setup(s.app)
}

func (s *Server) setupRoutes() {
	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	healthController := controller.NewHealthController(s.config, s.healthUsecase)
	authController := controller.NewRegistrationController(s.config, s.authUsecase)
	roleController := controller.NewRoleController(s.config, s.roleUsecase)
	userController := controller.NewUserController(s.config, s.userUsecase)

	s.setupHealthRoutes(v1, healthController)
	s.setupAuthRoutes(v1, authController)
	s.setupRoleRoutes(v1, roleController)
	s.setupUserRoutes(v1, userController)
}

func (s *Server) setupHealthRoutes(api fiber.Router, healthController *controller.HealthController) {
	health := api.Group("/health")
	health.Get("/", healthController.Check)
	health.Get("/ready", healthController.Ready)
	health.Get("/live", healthController.Live)
}

func (s *Server) setupAuthRoutes(api fiber.Router, registrationController *controller.AuthController) {
	auth := api.Group("/auth")

	auth.Post("/register", registrationController.Register)
	auth.Post("/register/special-account", registrationController.RegisterSpecialAccount)
	auth.Post("/login", registrationController.Login)
	auth.Post("/logout", registrationController.Logout)
	auth.Post("/verify-otp", registrationController.VerifyOTP)
	auth.Post("/complete-profile", registrationController.CompleteProfile)
	auth.Post("/resend-otp", registrationController.ResendOTP)
	auth.Post("/request-password-reset", registrationController.RequestPasswordReset)
	auth.Post("/reset-password", registrationController.ResetPassword)

	authProtected := auth.Group("")
	authProtected.Use(middleware.JWTAuth(s.config))
	authProtected.Post("/setup-pin", registrationController.SetupPIN)
}

func (s *Server) setupRoleRoutes(api fiber.Router, roleController *controller.RoleController) {
	roles := api.Group("/roles")

	roles.Use(middleware.JWTAuth(s.config))
	roles.Use(middleware.RequirePlatformAdmin())

	roles.Post("/", roleController.Create)
}

func (s *Server) setupUserRoutes(api fiber.Router, userController *controller.UserController) {
	users := api.Group("/users")
	users.Use(middleware.JWTAuth(s.config))

	// Self endpoints (any authenticated user)
	users.Get("/me", userController.GetMe)
	users.Put("/me", userController.UpdateMe)

	// Admin endpoints
	adminUsers := users.Group("")
	adminUsers.Use(middleware.RequirePlatformAdmin())

	adminUsers.Post("/", userController.Create)
	adminUsers.Get("/", userController.List)
	adminUsers.Get("/:id", userController.GetByID)
	adminUsers.Put("/:id", userController.Update)
	adminUsers.Delete("/:id", userController.Delete)
	adminUsers.Post("/:id/approve", userController.Approve)
	adminUsers.Post("/:id/reject", userController.Reject)
	adminUsers.Post("/:id/unlock", userController.Unlock)
	adminUsers.Post("/:id/reset-pin", userController.ResetPIN)
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting server on %s\n", addr)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func defaultErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
	})
}
