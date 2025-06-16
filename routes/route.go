package routes

import (
	"backend-attendance-deals/config"
	"backend-attendance-deals/middlewares"
	"backend-attendance-deals/repositories"
	"gorm.io/gorm"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type route struct {
	db              *gorm.DB
	routeController RouteController
}

func (r *route) SetupRoutes() *fiber.App {
	appEnv := config.GetEnv("APP_ENV", "DEVELOPMENT")
	log.Println(appEnv)
	userRepository := repositories.NewUserRepository(r.db)

	customMiddleware := middlewares.NewCustomMiddleware(
		userRepository,
	)
	app := fiber.New(fiber.Config{
		BodyLimit: 800 * 1024 * 1024, // this is the default limit of 4MB
	})

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Static("/file", "./public/upload")

	routeController := NewRouteController(r.db)

	app.Use(customMiddleware.FiberContextToContextMiddleware)

	app.Server().DisablePreParseMultipartForm = true
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Accept, Authorization, Content-Type, X-CSRF-Token, X-Request-ID",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
		AllowCredentials: false,
		MaxAge:           300,
	}))

	protectedApi := app.Group("/api", customMiddleware.AuthenticationMiddleware)
	protectedApi.Post("/login", routeController.authController.Login)

	protectedApi.Post("/attendance_period/create", routeController.attendancePeriodController.Create)
	protectedApi.Post("/attendance/create", routeController.attendanceController.Create)
	protectedApi.Post("/overtime/create", routeController.overtimeController.Create)
	protectedApi.Post("/reimbursement/create", routeController.reimbursementController.Create)
	protectedApi.Post("/payroll/create", routeController.payrollController.Create)
	protectedApi.Post("/payslip/create", routeController.payslipController.Create)
	protectedApi.Post("/payslip/summary", routeController.payslipController.List)

	return app
}

type RouteInterface interface {
	SetupRoutes() *fiber.App
}

func NewRoute(
	db *gorm.DB,
) RouteInterface {
	return &route{
		db: db,
	}
}
