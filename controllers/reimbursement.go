package controllers

import (
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/services"
	"github.com/gofiber/fiber/v2"
)

type ReimbursementControllerInterface interface {
	Create(c *fiber.Ctx) error
}

type reimbursementController struct {
	reimbursementService services.ReimbursementServiceInterface
}

func (a reimbursementController) Create(c *fiber.Ctx) error {
	var req shared.CreateReimbursementInput
	if err := c.BodyParser(&req); err != nil || req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	createReimbursement, err := a.reimbursementService.Create(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(createReimbursement)
}

func NewReimbursementController(
	reimbursementService services.ReimbursementServiceInterface,
) ReimbursementControllerInterface {
	return &reimbursementController{
		reimbursementService: reimbursementService,
	}
}
