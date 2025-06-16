package controllers

import (
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/services"
	"github.com/gofiber/fiber/v2"
)

type PaySlipControllerInterface interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
}

type payslipController struct {
	payslipService services.PayslipServiceInterface
}

func (a payslipController) List(c *fiber.Ctx) error {
	var req shared.ListSummaryPaySlip
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	createPayslip, err := a.payslipService.List(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(createPayslip)
}

func (a payslipController) Create(c *fiber.Ctx) error {
	var req shared.CreatePayslipInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	createPayslip, err := a.payslipService.Create(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(createPayslip)
}

func NewPaySlipController(
	payslipService services.PayslipServiceInterface,
) PaySlipControllerInterface {
	return &payslipController{
		payslipService: payslipService,
	}
}
