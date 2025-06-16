package routes

import (
	"backend-attendance-deals/controllers"
	"backend-attendance-deals/repositories"
	"backend-attendance-deals/services"
	"gorm.io/gorm"
)

type RouteController struct {
	authController             controllers.AuthControllerInterface
	attendancePeriodController controllers.AttendancePeriodControllerInterface
	attendanceController       controllers.AttendanceControllerInterface
	overtimeController         controllers.OvertimeControllerInterface
	reimbursementController    controllers.ReimbursementControllerInterface
	payrollController          controllers.PayrollControllerInterface
	payslipController          controllers.PaySlipControllerInterface
}

func NewRouteController(
	db *gorm.DB,
) RouteController {

	//Repository
	userRepository := repositories.NewUserRepository(db)
	auditLogRepository := repositories.NewLogAuditRepository(db)
	attendancePeriodRepository := repositories.NewAttendancePeriodRepository(db)
	attendanceRepository := repositories.NewAttendanceRepository(db)
	payrollRepository := repositories.NewPayrollRepository(db)
	overtimeRepository := repositories.NewOvertimeRepository(db)
	reimbursementRepository := repositories.NewReimbursementRepository(db)
	payslipRepository := repositories.NewPayslipRepository(db)

	//Service
	authService := services.NewAuthService(userRepository, auditLogRepository)
	attendancePeriodService := services.NewAttendancePeriodService(attendancePeriodRepository, auditLogRepository)
	attendanceService := services.NewAttendanceService(
		attendanceRepository,
		attendancePeriodRepository,
		auditLogRepository,
		payrollRepository,
	)
	overtimeService := services.NewOvertimeService(
		overtimeRepository,
		attendanceRepository,
		auditLogRepository,
	)
	reimbursementService := services.NewReimbursementService(
		reimbursementRepository, attendancePeriodRepository, auditLogRepository)
	payrollService := services.NewPayrollService(
		payrollRepository,
		attendanceRepository,
		overtimeRepository,
		reimbursementRepository,
		auditLogRepository,
		userRepository,
		attendancePeriodRepository,
		payslipRepository,
	)
	payslipService := services.NewPayslipService(
		payslipRepository,
		auditLogRepository,
		payrollRepository,
		userRepository,
		attendancePeriodRepository,
		attendanceRepository,
		overtimeRepository,
		reimbursementRepository,
	)

	//controller
	authController := controllers.NewAuthController(authService)
	attendancePeriodController := controllers.NewAttendancePeriodController(attendancePeriodService)
	attendanceController := controllers.NewAttendanceController(attendanceService)
	overtimeController := controllers.NewOvertimeController(overtimeService)
	reimbursementController := controllers.NewReimbursementController(reimbursementService)
	payrollController := controllers.NewPayrollController(payrollService)
	payslipController := controllers.NewPaySlipController(payslipService)

	return RouteController{
		authController:             authController,
		attendancePeriodController: attendancePeriodController,
		attendanceController:       attendanceController,
		overtimeController:         overtimeController,
		reimbursementController:    reimbursementController,
		payrollController:          payrollController,
		payslipController:          payslipController,
	}
}
