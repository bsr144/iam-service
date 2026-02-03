package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"
	"iam-service/delivery/http/router"
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

	// Infrastructure
	postgresDB, err := infrastructure.NewPostgres(cfg.Infra.Postgres, zapLogger)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}

	// Repositories
	authUserRepo := postgres.NewUserRepository(postgresDB)
	userProfileRepo := postgres.NewUserProfileRepository(postgresDB)
	userCredentialsRepo := postgres.NewUserCredentialsRepository(postgresDB)
	userSecurityRepo := postgres.NewUserSecurityRepository(postgresDB)
	emailVerificationRepo := postgres.NewEmailVerificationRepository(postgresDB)
	tenantRepo := postgres.NewTenantRepository(postgresDB)
	userActivationTrackingRepo := postgres.NewUserActivationTrackingRepository(postgresDB)
	roleRepo := postgres.NewRoleRepository(postgresDB)

	// External services
	emailService := mailer.NewEmailService()

	// Usecases
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

	// Controllers
	healthController := controller.NewHealthController(cfg, healthUsecase)
	authController := controller.NewRegistrationController(cfg, authUsecase)
	roleController := controller.NewRoleController(cfg, roleUsecase)
	userController := controller.NewUserController(cfg, userUsecase)

	server := &Server{
		app:    app,
		config: cfg,
		logger: zapLogger,
	}

	// Setup middleware
	mw := middleware.New(cfg, zapLogger)
	mw.Setup(app)

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	router.SetupHealthRoutes(v1, healthController)
	router.SetupAuthRoutes(v1, cfg, authController)
	router.SetupRoleRoutes(v1, cfg, roleController)
	router.SetupUserRoutes(v1, cfg, userController)

	return server
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
