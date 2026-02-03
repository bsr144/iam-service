package publickey

import (
	"encoding/base64"
	"iam-service/config"
	"iam-service/pkg/jwt"
	"math/big"
	"os"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Config *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		Config: cfg,
	}
}

func (h *Handler) GetPublicKeyPEM(c *fiber.Ctx) error {
	keyData, err := os.ReadFile(h.Config.JWT.PublicKeyPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read public key",
		})
	}

	c.Set("Content-Type", "application/x-pem-file")
	c.Set("Cache-Control", "public, max-age=3600")

	return c.Send(keyData)
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (h *Handler) GetJWKS(c *fiber.Ctx) error {

	publicKey, err := jwt.LoadPublicKeyFromFile(h.Config.JWT.PublicKeyPath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load public key",
		})
	}

	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Kid: "2024-01-key",
		N:   encodeBase64URL(publicKey.N.Bytes()),
		E:   encodeBase64URL(bigIntToBytes(publicKey.E)),
	}

	response := JWKSResponse{
		Keys: []JWK{jwk},
	}

	c.Set("Content-Type", "application/json")
	c.Set("Cache-Control", "public, max-age=3600")

	return c.JSON(response)
}

func encodeBase64URL(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func bigIntToBytes(n int) []byte {
	return big.NewInt(int64(n)).Bytes()
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
	wellKnown := router.Group("/.well-known")
	wellKnown.Get("/public-key.pem", h.GetPublicKeyPEM)
	wellKnown.Get("/jwks.json", h.GetJWKS)
}
