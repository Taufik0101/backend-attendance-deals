package middlewares

import (
	"backend-attendance-deals/config"
	"backend-attendance-deals/dto"
	"backend-attendance-deals/pkg/utils"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (m *CustomMiddleware) AuthenticationMiddleware(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization")

	ctx := c.UserContext()
	if len(authorizationHeader) > 0 && config.GetEnv("MIDDLEWARE_ENV", "DEVELOPMENT") != "PRODUCTION" {
		if !strings.Contains(authorizationHeader, "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "JWT Malformed"})
		}

		accessToken := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		accessTokenSecret := config.GetEnv("JWT_SECRET", "")
		validate, err := utils.JwtValidate(accessToken, accessTokenSecret)
		if err != nil || !validate.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		extractedClaims, err := utils.ExtractClaims(*validate)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		var accessTokenPayload dto.AccessTokenPayload
		err = mapstructure.Decode(extractedClaims["payload"], &accessTokenPayload)
		if err != nil {
			log.Println("Error decoding payload:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Payload malformed"})
		}

		ctx = context.WithValue(c.Context(), utils.CurrentUserKey, &accessTokenPayload)
		c.Locals(utils.CurrentUserKey, &accessTokenPayload)

	}

	// Request ID
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = uuid.NewString()
	}

	ctx = context.WithValue(c.Context(), utils.CurrentRequestKey, requestID)
	c.Locals(utils.CurrentRequestKey, requestID)

	// IP Address
	ctx = context.WithValue(c.Context(), utils.CurrentIPKey, c.IP())
	c.Locals(utils.CurrentIPKey, c.IP())

	c.SetUserContext(ctx)

	return c.Next()
}
