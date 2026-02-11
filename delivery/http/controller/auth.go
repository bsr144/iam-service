package controller

import (
	"strings"

	"iam-service/config"
	"iam-service/delivery/http/dto/response"
	"iam-service/delivery/http/presenter"
	"iam-service/iam/auth"
	"iam-service/iam/auth/authdto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func convertValidationErrors(errs validator.ValidationErrors) []errors.FieldError {
	result := make([]errors.FieldError, len(errs))
	for i, err := range errs {
		field := err.Field()
		var message string
		switch err.Tag() {
		case "required":
			message = field + " is required"
		case "email":
			message = field + " must be a valid email address"
		case "min":
			message = field + " must be at least " + err.Param() + " characters"
		case "max":
			message = field + " must be at most " + err.Param() + " characters"
		default:
			message = field + " is invalid"
		}
		result[i] = errors.FieldError{Field: field, Message: message}
	}
	return result
}

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
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.Register(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Registration initiated. Please check your email for OTP verification.",
		presenter.ToRegisterResponse(resp),
	))
}

func (rc *AuthController) RegisterSpecialAccount(c *fiber.Ctx) error {
	var req authdto.RegisterSpecialAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.RegisterSpecialAccount(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"Special Account Registration is successful.",
		presenter.ToRegisterSpecialAccountResponse(resp),
	))
}

func (rc *AuthController) Login(c *fiber.Ctx) error {
	var req authdto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.Login(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Login successful",
		presenter.ToLoginResponse(resp),
	))
}

func (rc *AuthController) VerifyOTP(c *fiber.Ctx) error {
	var req authdto.VerifyOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.VerifyOTP(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"OTP verified successfully.",
		presenter.ToVerifyOTPResponse(resp),
	))
}

func (rc *AuthController) CompleteProfile(c *fiber.Ctx) error {
	var req authdto.CompleteProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.CompleteProfile(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToCompleteProfileResponse(resp),
	))
}

func (rc *AuthController) ResendOTP(c *fiber.Ctx) error {
	var req authdto.ResendOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.ResendOTP(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"OTP resent successfully. Please check your email.",
		presenter.ToResendOTPResponse(resp),
	))
}

func (rc *AuthController) Logout(c *fiber.Ctx) error {
	var req authdto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	err := rc.authUsecase.Logout(c.UserContext(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Logout successful",
		nil,
	))
}

func (rc *AuthController) SetupPIN(c *fiber.Ctx) error {
	var req authdto.SetupPINRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	resp, err := rc.authUsecase.SetupPIN(c.Context(), userID, &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"PIN setup successful",
		presenter.ToSetupPINResponse(resp),
	))
}

func (rc *AuthController) RequestPasswordReset(c *fiber.Ctx) error {
	var req authdto.RequestPasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.RequestPasswordReset(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Password reset OTP has been sent to your email",
		presenter.ToRequestPasswordResetResponse(resp),
	))
}

func (rc *AuthController) ResetPassword(c *fiber.Ctx) error {
	var req authdto.ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.ResetPassword(c.Context(), &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToResetPasswordResponse(resp),
	))
}

func (rc *AuthController) InitiateRegistration(c *fiber.Ctx) error {
	tenantID, err := getTenantIDFromHeader(c)
	if err != nil {
		return err
	}

	var req authdto.InitiateRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	ipAddress := getClientIP(c).String()
	userAgent := getUserAgent(c)

	resp, err := rc.authUsecase.InitiateRegistration(c.Context(), tenantID, &req, ipAddress, userAgent)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToInitiateRegistrationResponse(resp),
	))
}

func (rc *AuthController) VerifyRegistrationOTP(c *fiber.Ctx) error {
	tenantID, err := getTenantIDFromHeader(c)
	if err != nil {
		return err
	}

	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	var req authdto.VerifyRegistrationOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.VerifyRegistrationOTP(c.Context(), tenantID, registrationID, &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToVerifyRegistrationOTPResponse(resp),
	))
}

func (rc *AuthController) ResendRegistrationOTP(c *fiber.Ctx) error {
	tenantID, err := getTenantIDFromHeader(c)
	if err != nil {
		return err
	}

	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	var req authdto.ResendRegistrationOTPRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	resp, err := rc.authUsecase.ResendRegistrationOTP(c.Context(), tenantID, registrationID, &req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToResendRegistrationOTPResponse(resp),
	))
}

func (rc *AuthController) GetRegistrationStatus(c *fiber.Ctx) error {
	tenantID, err := getTenantIDFromHeader(c)
	if err != nil {
		return err
	}

	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	email := c.Query("email")
	if email == "" {
		return errors.ErrBadRequest("email query parameter is required")
	}

	resp, err := rc.authUsecase.GetRegistrationStatus(c.Context(), tenantID, registrationID, email)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Registration status retrieved",
		presenter.ToRegistrationStatusResponse(resp),
	))
}

func (rc *AuthController) CompleteRegistration(c *fiber.Ctx) error {
	tenantID, err := getTenantIDFromHeader(c)
	if err != nil {
		return err
	}

	registrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errors.ErrBadRequest("Invalid registration ID format")
	}

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return errors.ErrUnauthorized("Authorization header is required")
	}

	registrationToken := strings.TrimPrefix(authHeader, "Bearer ")
	if registrationToken == authHeader {
		return errors.ErrUnauthorized("Invalid authorization format. Use: Bearer <token>")
	}

	var req authdto.CompleteRegistrationRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.ErrBadRequest("Invalid request body")
	}

	if err := rc.validate.Struct(&req); err != nil {
		return errors.ErrValidationWithFields(convertValidationErrors(err.(validator.ValidationErrors)))
	}

	ipAddress := getClientIP(c).String()
	userAgent := getUserAgent(c)

	resp, err := rc.authUsecase.CompleteRegistration(c.Context(), tenantID, registrationID, registrationToken, &req, ipAddress, userAgent)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		resp.Message,
		presenter.ToCompleteRegistrationResponse(resp),
	))
}
