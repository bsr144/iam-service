package router

import (
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupParticipantRoutes(api fiber.Router, ctrl *controller.ParticipantController, jwtMiddleware fiber.Handler) {
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
