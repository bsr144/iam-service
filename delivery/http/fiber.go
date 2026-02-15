package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iam-service/config"
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/dto/response"
	"iam-service/delivery/http/middleware"
	"iam-service/delivery/http/router"
	"iam-service/iam/auth"
	"iam-service/iam/health"
	"iam-service/iam/role"
	"iam-service/iam/user"
	"iam-service/impl/mailer"
	"iam-service/impl/postgres"
	implredis "iam-service/impl/redis"
	"iam-service/infrastructure"
	"iam-service/masterdata"
	apperrors "iam-service/pkg/errors"
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
	zapLogger, _ := logger.NewZapLoggerWithConfig(cfg.Log, cfg.App.Environment)
	auditLogger := logger.NewAuditLogger(zapLogger, logger.AuditConfig{
		Enabled: cfg.Log.AuditEnabled,
	})

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		AppName:      cfg.App.Name,
		BodyLimit:    256 * 1024,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: createErrorHandler(cfg, zapLogger),
	})

	postgresDB, err := infrastructure.NewPostgres(cfg.Infra.Postgres, zapLogger)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}

	redisClient, err := infrastructure.NewRedis(cfg.Infra.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis:", err)
	}
	redisWrapper := implredis.NewRedis(redisClient)

	txManager := postgres.NewTransactionManager(postgresDB)

	authUserRepo := postgres.NewUserRepository(postgresDB)
	userProfileRepo := postgres.NewUserProfileRepository(postgresDB)
	userAuthMethodRepo := postgres.NewUserAuthMethodRepository(postgresDB)
	userSecurityStateRepo := postgres.NewUserSecurityStateRepository(postgresDB)
	tenantRepo := postgres.NewTenantRepository(postgresDB)
	roleRepo := postgres.NewRoleRepository(postgresDB)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(postgresDB)
	userRoleRepo := postgres.NewUserRoleRepository(postgresDB)
	productRepo := postgres.NewProductRepository(postgresDB)
	permissionRepo := postgres.NewPermissionRepository(postgresDB)
	rolePermissionRepo := postgres.NewRolePermissionRepository(postgresDB)
	userSessionRepo := postgres.NewUserSessionRepository(postgresDB)
	userTenantRegRepo := postgres.NewUserTenantRegistrationRepository(postgresDB)
	productsByTenantRepo := postgres.NewProductsByTenantRepository(postgresDB)

	masterdataCategoryRepo := postgres.NewMasterdataCategoryRepository(postgresDB)
	masterdataItemRepo := postgres.NewMasterdataItemRepository(postgresDB)

	emailService := mailer.NewEmailService(&cfg.Email)

	healthUsecase := health.NewUsecase()
	authUsecase := auth.NewUsecase(
		txManager,
		cfg,
		authUserRepo,
		userProfileRepo,
		userAuthMethodRepo,
		userSecurityStateRepo,
		tenantRepo,
		roleRepo,
		refreshTokenRepo,
		userRoleRepo,
		productRepo,
		permissionRepo,
		emailService,
		redisWrapper,
		redisWrapper,
		userSessionRepo,
		userTenantRegRepo,
		productsByTenantRepo,
		auditLogger,
	)
	roleUsecase := role.NewUsecase(
		txManager,
		cfg,
		tenantRepo,
		roleRepo,
		rolePermissionRepo,
	)
	userUsecase := user.NewUsecase(
		txManager,
		cfg,
		authUserRepo,
		userProfileRepo,
		userAuthMethodRepo,
		userSecurityStateRepo,
		tenantRepo,
		roleRepo,
		userRoleRepo,
	)
	masterdataUsecase := masterdata.NewUsecase(
		cfg,
		masterdataCategoryRepo,
		masterdataItemRepo,
		redisWrapper,
	)

	healthController := controller.NewHealthController(cfg, healthUsecase)
	authController := controller.NewRegistrationController(cfg, authUsecase)
	roleController := controller.NewRoleController(cfg, roleUsecase)
	userController := controller.NewUserController(cfg, userUsecase)
	masterdataController := controller.NewMasterdataController(cfg, masterdataUsecase)

	server := &Server{
		app:    app,
		config: cfg,
		logger: zapLogger,
	}

	mw := middleware.New(cfg, zapLogger)
	mw.Setup(app)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	router.SetupHealthRoutes(v1, healthController)

	iam := v1.Group("/iam")
	router.SetupAuthRoutes(iam, cfg, authController)
	router.SetupRoleRoutes(iam, cfg, roleController)
	router.SetupUserRoutes(iam, cfg, userController)

	router.SetupMasterdataRoutes(v1, cfg, masterdataController)

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

func createErrorHandler(cfg *config.Config, zapLogger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		requestID := middleware.GetRequestID(c)
		includeDebug := cfg.IsDevelopment()

		var appErr *apperrors.AppError
		if apperrors.As(err, &appErr) {

			logFields := []zap.Field{
				zap.String("request_id", requestID),
				zap.String("code", appErr.Code),
				zap.String("message", appErr.Message),
				zap.String("kind", appErr.Kind.String()),
				zap.String("file", appErr.File),
				zap.Int("line", appErr.Line),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			}
			if appErr.Op != "" {
				logFields = append(logFields, zap.String("op", appErr.Op))
			}
			if appErr.Err != nil {
				logFields = append(logFields, zap.Error(appErr.Err))
			}

			if appErr.HTTPStatus >= 500 {
				zapLogger.Error(appErr.Message, logFields...)
			} else if appErr.HTTPStatus >= 400 {
				zapLogger.Warn(appErr.Message, logFields...)
			}

			resp := response.APIResponse{
				Success:   false,
				Error:     appErr.Code,
				Message:   appErr.Message,
				RequestID: requestID,
			}

			if appErr.Code == apperrors.CodeValidation && appErr.Details != nil {
				if fieldErrors, ok := appErr.Details["fields"].([]apperrors.FieldError); ok {
					resp.Errors = make([]response.FieldError, len(fieldErrors))
					for i, fe := range fieldErrors {
						resp.Errors[i] = response.FieldError{
							Field:   fe.Field,
							Message: fe.Message,
						}
					}
				}
			}

			if includeDebug && appErr.Err != nil {
				resp.Debug = &response.DebugInfo{
					Cause: appErr.Err.Error(),
				}
			}

			return c.Status(appErr.HTTPStatus).JSON(resp)
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			zapLogger.Warn("Fiber error",
				zap.String("request_id", requestID),
				zap.Int("status", fiberErr.Code),
				zap.String("message", fiberErr.Message),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			)
			return c.Status(fiberErr.Code).JSON(response.APIResponse{
				Success:   false,
				Error:     "FIBER_ERROR",
				Message:   fiberErr.Message,
				RequestID: requestID,
			})
		}

		zapLogger.Error("Unexpected error",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
		)

		resp := response.APIResponse{
			Success:   false,
			Error:     "INTERNAL_SERVER_ERROR",
			Message:   "an unexpected error occurred",
			RequestID: requestID,
		}

		if includeDebug {
			resp.Debug = &response.DebugInfo{
				Cause: err.Error(),
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}
}
