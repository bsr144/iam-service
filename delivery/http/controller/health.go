package controller

import (
	"iam-service/config"
	"iam-service/internal/health"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	config        *config.Config
	healthUsecase health.Usecase
}

func NewHealthController(cfg *config.Config, healthUsecase health.Usecase) *HealthController {
	return &HealthController{
		config:        cfg,
		healthUsecase: healthUsecase,
	}
}

func (h *HealthController) Check(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"status":      "healthy",
			"app":         h.config.App.Name,
			"version":     h.config.App.Version,
			"environment": h.config.App.Environment,
		},
	})
}

func (h *HealthController) Ready(c *fiber.Ctx) error {
	if err := h.healthUsecase.CheckHealth(); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"success": false,
			"data": fiber.Map{
				"status": "not ready",
			},
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"status": "ready",
		},
	})
}

func (h *HealthController) Live(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"status": "alive",
		},
	})
}
