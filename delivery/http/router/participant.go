package router

import (
	"time"

	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func selfRegRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if uid, ok := c.Locals("userID").(string); ok && uid != "" {
				return "self-reg:" + uid
			}
			return "self-reg:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many requests",
				"code":    "ERR_TOO_MANY_REQUESTS",
			})
		},
	})
}

func SetupParticipantRoutes(api fiber.Router, ctrl *controller.ParticipantController, jwtMiddleware fiber.Handler) {
	selfReg := api.Group("/participants")
	selfReg.Use(jwtMiddleware)
	selfReg.Post("/self-register", selfRegRateLimit(), ctrl.SelfRegister)

	participants := api.Group("/products/:productId/participants")
	participants.Use(jwtMiddleware)
	participants.Use(middleware.ExtractTenantContext())
	participants.Use(middleware.ExtractProductContext())

	participants.Post("/",
		middleware.RequireProductPermission("participant:create"),
		ctrl.Create,
	)

	participants.Get("/",
		middleware.RequireProductPermission("participant:read"),
		ctrl.List,
	)

	participants.Get("/:id",
		middleware.RequireProductPermission("participant:read"),
		ctrl.Get,
	)

	participants.Put("/:id/personal-data",
		middleware.RequireProductPermission("participant:update"),
		ctrl.UpdatePersonalData,
	)

	participants.Put("/:id/identities",
		middleware.RequireProductPermission("participant:update"),
		ctrl.SaveIdentity,
	)
	participants.Delete("/:id/identities/:identityId",
		middleware.RequireProductPermission("participant:update"),
		ctrl.DeleteIdentity,
	)

	participants.Put("/:id/addresses",
		middleware.RequireProductPermission("participant:update"),
		ctrl.SaveAddress,
	)
	participants.Delete("/:id/addresses/:addressId",
		middleware.RequireProductPermission("participant:update"),
		ctrl.DeleteAddress,
	)

	participants.Put("/:id/bank-accounts",
		middleware.RequireProductPermission("participant:update"),
		ctrl.SaveBankAccount,
	)
	participants.Delete("/:id/bank-accounts/:accountId",
		middleware.RequireProductPermission("participant:update"),
		ctrl.DeleteBankAccount,
	)

	participants.Put("/:id/family-members",
		middleware.RequireProductPermission("participant:update"),
		ctrl.SaveFamilyMember,
	)
	participants.Delete("/:id/family-members/:memberId",
		middleware.RequireProductPermission("participant:update"),
		ctrl.DeleteFamilyMember,
	)

	participants.Put("/:id/employment",
		middleware.RequireProductPermission("participant:update"),
		ctrl.SaveEmployment,
	)

	participants.Put("/:id/pension",
		middleware.RequireProductPermission("participant:update"),
		ctrl.SavePension,
	)

	participants.Put("/:id/beneficiaries",
		middleware.RequireProductPermission("participant:update"),
		ctrl.SaveBeneficiary,
	)
	participants.Delete("/:id/beneficiaries/:beneficiaryId",
		middleware.RequireProductPermission("participant:update"),
		ctrl.DeleteBeneficiary,
	)

	participants.Post("/:id/files",
		middleware.RequireProductPermission("participant:update"),
		ctrl.UploadFile,
	)

	participants.Get("/:id/status-history",
		middleware.RequireProductPermission("participant:read"),
		ctrl.GetStatusHistory,
	)

	participants.Post("/:id/submit",
		middleware.RequireProductPermission("participant:submit"),
		ctrl.Submit,
	)

	participants.Post("/:id/approve",
		middleware.RequireProductPermission("participant:approve"),
		ctrl.Approve,
	)

	participants.Post("/:id/reject",
		middleware.RequireProductPermission("participant:reject"),
		ctrl.Reject,
	)

	participants.Delete("/:id",
		middleware.RequireProductPermission("participant:delete"),
		ctrl.Delete,
	)
}
