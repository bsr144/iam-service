package router

import (
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupParticipantRoutes(api fiber.Router, ctrl *controller.ParticipantController, jwtMiddleware fiber.Handler) {
	participants := api.Group("/participants")
	participants.Use(jwtMiddleware)
	participants.Use(middleware.ExtractTenantContext())

	participants.Post("/",
		middleware.RequireTenantPermission("participant:create"),
		ctrl.Create,
	)

	participants.Get("/",
		middleware.RequireTenantPermission("participant:read"),
		ctrl.List,
	)

	participants.Get("/:id",
		middleware.RequireTenantPermission("participant:read"),
		ctrl.Get,
	)

	participants.Put("/:id/personal-data",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.UpdatePersonalData,
	)

	// Identity sub-resource
	participants.Put("/:id/identities",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.SaveIdentity,
	)
	participants.Delete("/:id/identities/:identityId",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.DeleteIdentity,
	)

	// Address sub-resource
	participants.Put("/:id/addresses",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.SaveAddress,
	)
	participants.Delete("/:id/addresses/:addressId",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.DeleteAddress,
	)

	// Bank account sub-resource
	participants.Put("/:id/bank-accounts",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.SaveBankAccount,
	)
	participants.Delete("/:id/bank-accounts/:accountId",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.DeleteBankAccount,
	)

	// Family member sub-resource
	participants.Put("/:id/family-members",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.SaveFamilyMember,
	)
	participants.Delete("/:id/family-members/:memberId",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.DeleteFamilyMember,
	)

	// Employment sub-resource
	participants.Put("/:id/employment",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.SaveEmployment,
	)

	// Beneficiary sub-resource
	participants.Put("/:id/beneficiaries",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.SaveBeneficiary,
	)
	participants.Delete("/:id/beneficiaries/:beneficiaryId",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.DeleteBeneficiary,
	)

	// File upload (5MB limit for this route)
	participants.Post("/:id/files",
		middleware.RequireTenantPermission("participant:update"),
		ctrl.UploadFile,
	)

	// Status history
	participants.Get("/:id/status-history",
		middleware.RequireTenantPermission("participant:read"),
		ctrl.GetStatusHistory,
	)

	// Workflow actions
	participants.Post("/:id/submit",
		middleware.RequireTenantPermission("participant:submit"),
		ctrl.Submit,
	)

	participants.Post("/:id/approve",
		middleware.RequireTenantPermission("participant:approve"),
		ctrl.Approve,
	)

	participants.Post("/:id/reject",
		middleware.RequireTenantPermission("participant:reject"),
		ctrl.Reject,
	)

	participants.Delete("/:id",
		middleware.RequireTenantPermission("participant:delete"),
		ctrl.Delete,
	)
}
