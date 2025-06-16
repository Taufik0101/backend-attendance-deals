package controllers

import (
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/services"
	"github.com/gofiber/fiber/v2"
)

type AuthControllerInterface interface {
	Login(c *fiber.Ctx) error
}

type authController struct {
	authService services.AuthServiceInterface
}

func (a authController) Login(c *fiber.Ctx) error {
	var req shared.LoginInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	login, err := a.authService.Login(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(login)
}

func NewAuthController(
	authService services.AuthServiceInterface,
) AuthControllerInterface {
	return &authController{
		authService: authService,
	}
}
