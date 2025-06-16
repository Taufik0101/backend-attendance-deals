package middlewares

import (
	"backend-attendance-deals/repositories"
	"github.com/gofiber/fiber/v2"
)

type CustomMiddleware struct {
	UserRepository repositories.UserRepositoryInterface
}

type CustomMiddlewareInterface interface {
	AuthenticationMiddleware(ctx *fiber.Ctx) error
	FiberContextToContextMiddleware(ctx *fiber.Ctx) error
}

func NewCustomMiddleware(
	userRepository repositories.UserRepositoryInterface,
) CustomMiddlewareInterface {
	return &CustomMiddleware{
		UserRepository: userRepository,
	}
}
