package controller

import (
	"iam-service/config"
	"iam-service/internal/user"
	"iam-service/internal/user/userdto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	config      *config.Config
	userUsecase user.Usecase
	validate    *validator.Validate
}

func NewUserController(cfg *config.Config, userUsecase user.Usecase) *UserController {
	return &UserController{
		config:      cfg,
		userUsecase: userUsecase,
		validate:    validate,
	}
}

func (uc *UserController) Create(c *fiber.Ctx) error {
	var req userdto.CreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := uc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := uc.userUsecase.Create(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse(
		"User created successfully",
		resp,
	))
}
