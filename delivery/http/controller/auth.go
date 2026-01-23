package controller

import (
	"iam-service/config"
	"iam-service/delivery/http/middleware"
	"iam-service/internal/auth"
	"iam-service/internal/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	config      *config.Config
	authUsecase auth.Usecase
	validate    *validator.Validate
}

func NewRegistrationController(cfg *config.Config, authUsecase auth.Usecase) *AuthController {
	return &AuthController{
		config:      cfg,
		authUsecase: authUsecase,
		validate:    validate,
	}
}

func (rc *AuthController) Register(c *fiber.Ctx) error {
	var req authdto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.Register(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse(
		"Registration initiated. Please check your email for OTP verification.",
		resp,
	))
}

func (rc *AuthController) RegisterSpecialAccount(c *fiber.Ctx) error {
	var req authdto.RegisterSpecialAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.RegisterSpecialAccount(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse(
		"Special Account Registration is successful.",
		resp,
	))
}

func (rc *AuthController) Login(c *fiber.Ctx) error {
	var req authdto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.Login(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		"Login successful",
		resp,
	))
}

func (rc *AuthController) VerifyOTP(c *fiber.Ctx) error {
	var req authdto.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.VerifyOTP(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		"OTP verified successfully.",
		resp,
	))
}

func (rc *AuthController) CompleteProfile(c *fiber.Ctx) error {
	var req authdto.CompleteProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.CompleteProfile(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		resp.Message,
		resp,
	))
}

func (rc *AuthController) ResendOTP(c *fiber.Ctx) error {
	var req authdto.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.ResendOTP(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		"OTP resent successfully. Please check your email.",
		resp,
	))
}

func (rc *AuthController) Logout(c *fiber.Ctx) error {
	var req authdto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	err := rc.authUsecase.Logout(req.RefreshToken)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		"Logout successful",
		nil,
	))
}

func (rc *AuthController) SetupPIN(c *fiber.Ctx) error {
	var req authdto.SetupPINRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	resp, err := rc.authUsecase.SetupPIN(c.Context(), userID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		"PIN setup successful",
		resp,
	))
}

func (rc *AuthController) RequestPasswordReset(c *fiber.Ctx) error {
	var req authdto.RequestPasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.RequestPasswordReset(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		"Password reset OTP has been sent to your email",
		resp,
	))
}

func (rc *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req authdto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := rc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := rc.authUsecase.ResetPassword(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse(
		resp.Message,
		resp,
	))
}

func handleError(c *fiber.Ctx, err error) error {
	appErr := errors.GetAppError(err)
	if appErr != nil {
		return c.Status(appErr.HTTPStatus).JSON(ErrorResponse(
			appErr.Code,
			appErr.Message,
		))
	}

	return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse(
		errors.CodeInternal,
		"An unexpected error occurred",
	))
}

func formatValidationErrors(errs validator.ValidationErrors) map[string]string {
	result := make(map[string]string)
	for _, err := range errs {
		field := err.Field()
		switch err.Tag() {
		case "required":
			result[field] = field + " is required"
		case "email":
			result[field] = field + " must be a valid email address"
		case "min":
			result[field] = field + " must be at least " + err.Param() + " characters"
		case "max":
			result[field] = field + " must be at most " + err.Param() + " characters"
		default:
			result[field] = field + " is invalid"
		}
	}
	return result
}
