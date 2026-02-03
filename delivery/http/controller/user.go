package controller

import (
	"iam-service/config"
	"iam-service/delivery/http/dto/response"
	"iam-service/delivery/http/middleware"
	"iam-service/iam/user"
	"iam-service/iam/user/userdto"
	"iam-service/pkg/errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := uc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := uc.userUsecase.Create(c.Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse(
		"User created successfully",
		resp,
	))
}

func (uc *UserController) GetMe(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	resp, err := uc.userUsecase.GetMe(c.Context(), userID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"User profile retrieved successfully",
		resp,
	))
}

func (uc *UserController) UpdateMe(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	var req userdto.UpdateMeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := uc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := uc.userUsecase.UpdateMe(c.Context(), userID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"Profile updated successfully",
		resp,
	))
}

func (uc *UserController) List(c *fiber.Ctx) error {
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		return handleError(c, err)
	}

	var req userdto.ListRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid query parameters",
		))
	}

	resp, err := uc.userUsecase.List(c.Context(), tenantID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.APIResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    resp.Users,
		Pagination: &response.Pagination{
			Total:      resp.Pagination.Total,
			Page:       resp.Pagination.Page,
			Limit:      resp.Pagination.PerPage,
			TotalPages: resp.Pagination.TotalPages,
		},
	})
}

func (uc *UserController) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid user ID",
		))
	}

	resp, err := uc.userUsecase.GetByID(c.Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"User retrieved successfully",
		resp,
	))
}

func (uc *UserController) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid user ID",
		))
	}

	var req userdto.UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := uc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := uc.userUsecase.Update(c.Context(), id, &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"User updated successfully",
		resp,
	))
}

func (uc *UserController) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid user ID",
		))
	}

	if err := uc.userUsecase.Delete(c.Context(), id); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		"User deleted successfully",
		nil,
	))
}

func (uc *UserController) Approve(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid user ID",
		))
	}

	approverID, err := middleware.GetUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	resp, err := uc.userUsecase.Approve(c.Context(), id, approverID)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		resp,
	))
}

func (uc *UserController) Reject(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid user ID",
		))
	}

	approverID, err := middleware.GetUserID(c)
	if err != nil {
		return handleError(c, err)
	}

	var req userdto.RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid request body",
		))
	}

	if err := uc.validate.Struct(&req); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponseWithDetails(
			errors.CodeValidation,
			"Validation failed",
			formatValidationErrors(validationErrors),
		))
	}

	resp, err := uc.userUsecase.Reject(c.Context(), id, approverID, &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		resp,
	))
}

func (uc *UserController) Unlock(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid user ID",
		))
	}

	resp, err := uc.userUsecase.Unlock(c.Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		resp,
	))
}

func (uc *UserController) ResetPIN(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(
			errors.CodeBadRequest,
			"Invalid user ID",
		))
	}

	resp, err := uc.userUsecase.ResetPIN(c.Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse(
		resp.Message,
		resp,
	))
}
