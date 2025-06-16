package controllers

import (
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/services"
	"github.com/gofiber/fiber/v2"
)

type AttendanceControllerInterface interface {
	Create(c *fiber.Ctx) error
}

type attendanceController struct {
	attendanceService services.AttendanceServiceInterface
}

func (a attendanceController) Create(c *fiber.Ctx) error {
	var req shared.CreateAttendanceInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	createAttendance, err := a.attendanceService.Create(c.UserContext(), req)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(createAttendance)
}

func NewAttendanceController(
	attendanceService services.AttendanceServiceInterface,
) AttendanceControllerInterface {
	return &attendanceController{
		attendanceService: attendanceService,
	}
}
