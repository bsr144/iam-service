package middleware

import (
	"net"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ClientIPKey        = "client_ip"
	UserAgentKey       = "user_agent"
	TenantIDFromHdrKey = "tenant_id_from_header"
)

func RequestContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientIP := extractClientIP(c)
		c.Locals(ClientIPKey, clientIP)

		userAgent := c.Get("User-Agent")
		c.Locals(UserAgentKey, userAgent)

		if tenantIDStr := c.Get("X-Tenant-ID"); tenantIDStr != "" {
			if tenantID, err := uuid.Parse(tenantIDStr); err == nil {
				c.Locals(TenantIDFromHdrKey, tenantID)
			}
		}

		return c.Next()
	}
}

func extractClientIP(c *fiber.Ctx) net.IP {

	if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
		return net.ParseIP(forwarded)
	}

	return net.ParseIP(c.IP())
}
