package controllers

import (
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/services"
	"github.com/gofiber/fiber/v2"
)

type PayrollControllerInterface interface {
	Create(c *fiber.Ctx) error
}

type payrollController struct {
	payrollService services.PayrollServiceInterface
}

func (a payrollController) Create(c *fiber.Ctx) error {
	var req shared.CreatePayrollInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	createPayroll, err := a.payrollService.Create(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(createPayroll)
}

func NewPayrollController(
	payrollService services.PayrollServiceInterface,
) PayrollControllerInterface {
	return &payrollController{
		payrollService: payrollService,
	}
}
