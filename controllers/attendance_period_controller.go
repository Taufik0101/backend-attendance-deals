package controllers

import (
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/services"
	"github.com/gofiber/fiber/v2"
)

type AttendancePeriodControllerInterface interface {
	Create(c *fiber.Ctx) error
}

type attendancePeriodController struct {
	attendancePeriodService services.AttendancePeriodServiceInterface
}

func (a attendancePeriodController) Create(c *fiber.Ctx) error {
	var req shared.CreateAttendancePeriodInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createAttendancePeriod, err := a.attendancePeriodService.Create(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(createAttendancePeriod)
}

func NewAttendancePeriodController(
	attendancePeriodService services.AttendancePeriodServiceInterface,
) AttendancePeriodControllerInterface {
	return &attendancePeriodController{
		attendancePeriodService: attendancePeriodService,
	}
}
