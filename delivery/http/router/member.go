package router

import (
	"iam-service/delivery/http/controller"
	"iam-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupMemberRoutes(api fiber.Router, ctrl *controller.MemberController, jwtMiddleware fiber.Handler) {
	products := api.Group("/products/:productId/members")
	products.Use(jwtMiddleware)
	products.Use(middleware.ExtractTenantContext())
	products.Use(middleware.ExtractProductContext())

	products.Post("/register", ctrl.Register)

	adminMiddleware := middleware.RequireProductRole("TENANT_PRODUCT_ADMIN")

	products.Get("/", adminMiddleware, ctrl.List)
	products.Get("/:memberId", adminMiddleware, ctrl.Get)
	products.Post("/:memberId/approve", adminMiddleware, ctrl.Approve)
	products.Post("/:memberId/reject", adminMiddleware, ctrl.Reject)
	products.Put("/:memberId/role", adminMiddleware, ctrl.ChangeRole)
	products.Post("/:memberId/deactivate", adminMiddleware, ctrl.Deactivate)
}
