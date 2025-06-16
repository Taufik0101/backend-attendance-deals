package controllers

import (
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/services"
	"github.com/gofiber/fiber/v2"
)

type OvertimeControllerInterface interface {
	Create(c *fiber.Ctx) error
}

type overtimeController struct {
	overtimeService services.OvertimeServiceInterface
}

func (a overtimeController) Create(c *fiber.Ctx) error {
	var req shared.CreateOvertimeInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	createOvertime, err := a.overtimeService.Create(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(createOvertime)
}

func NewOvertimeController(
	overtimeService services.OvertimeServiceInterface,
) OvertimeControllerInterface {
	return &overtimeController{
		overtimeService: overtimeService,
	}
}
