package controller

import (
	"iam-service/config"
	"iam-service/internal/role"
	"iam-service/internal/role/roledto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type RoleController struct {
	config      *config.Config
	roleUsecase role.Usecase
	validate    *validator.Validate
}

func NewRoleController(cfg *config.Config, roleUsecase role.Usecase) *RoleController {
	return &RoleController{
		config:      cfg,
		roleUsecase: roleUsecase,
		validate:    validate,
	}
}

func (rc *RoleController) Create(c *fiber.Ctx) error {
	var req roledto.CreateRequest
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

	resp, err := rc.roleUsecase.Create(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponse(
		"Role created successfully",
		resp,
	))
}
