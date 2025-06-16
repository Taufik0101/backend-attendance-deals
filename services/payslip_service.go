package services

import (
	"backend-attendance-deals/config"
	"backend-attendance-deals/database/models"
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/pkg/utils"
	"backend-attendance-deals/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
	"sort"
	"strings"
	"time"
)

type PayslipServiceInterface interface {
	Create(ctx context.Context, input shared.CreatePayslipInput) (*shared.CreatePayslipResponse, error)
	List(ctx context.Context, input shared.ListSummaryPaySlip) (*shared.PaySlipResponse, error)
}

type payslipService struct {
	payslipRepository       repositories.PayslipRepositoryInterface
	logAuditRepository      repositories.LogAuditRepositoryInterface
	payrollRepository       repositories.PayrollRepositoryInterface
	userRepository          repositories.UserRepositoryInterface
	periodRepository        repositories.AttendancePeriodRepositoryInterface
	attendanceRepository    repositories.AttendanceRepositoryInterface
	overtimeRepository      repositories.OvertimeRepositoryInterface
	reimbursementRepository repositories.ReimbursementRepositoryInterface
}

func (a payslipService) List(ctx context.Context, input shared.ListSummaryPaySlip) (*shared.PaySlipResponse, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)
	currentUser := utils.GetCurrentUser(&ctx)
	UUIDUser, _ := uuid.Parse(currentUser.UserId)
	UUIDPeriod, _ := uuid.Parse(input.PeriodID)

	if strings.ToLower(currentUser.UserType) != "admin" {
		return nil, errors.New("FORBIDDEN")
	}

	output := new(shared.PaySlipResponse)

	arrDetailPaySlip := make([]shared.DetailPaySlipResponse, 0)
	totalTakeHomePay := float64(0)

	findPayroll, _ := a.payrollRepository.FindOneWithJoin(map[string]any{
		"period_id": UUIDPeriod,
	}, nil, nil, nil, nil, nil)

	if findPayroll == nil {
		return nil, errors.New("period id not found")
	}

	findPayslip, _ := a.payslipRepository.FindWithJoin(
		map[string]any{
			"payroll_id": findPayroll.ID,
		},
		nil,
		nil,
		nil,
		nil,
		[]string{"User"},
	)

	for _, ps := range findPayslip {
		user := ""
		if ps.User != nil {
			user = ps.User.Username
		}
		arrDetailPaySlip = append(arrDetailPaySlip, shared.DetailPaySlipResponse{
			User:        user,
			TakeHomePay: ps.TotalPay,
		})
		totalTakeHomePay += ps.TotalPay
	}

	output.TotalTakeHomePay = totalTakeHomePay
	output.DetailTakeHomePay = arrDetailPaySlip

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		UserID:    UUIDUser,
		Action:    "get_summary_payslip",
		Entity:    "PaySlip",
		EntityID:  UUIDUser,
		IpAddress: *ipAddress,
		RequestID: *requestID,
		PeriodID:  &UUIDPeriod,
	})

	_, err := a.logAuditRepository.Create(ctx, arrLogAudit)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return output, nil
}

func (a payslipService) Create(ctx context.Context, input shared.CreatePayslipInput) (*shared.CreatePayslipResponse, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)
	currentUser := utils.GetCurrentUser(&ctx)
	UUIDUser, _ := uuid.Parse(currentUser.UserId)
	UUIDPeriod, _ := uuid.Parse(input.PeriodID)

	if strings.ToLower(currentUser.UserType) != "employee" {
		return nil, errors.New("FORBIDDEN")
	}

	output := new(shared.CreatePayslipResponse)

	// findUser
	findUser, err := a.userRepository.FindOneWithAttribute(
		map[string]any{
			"id": UUIDUser,
		},
		nil,
		nil,
	)

	if err != nil {
		return nil, err
	}

	// find period
	findAttendancePeriod, _ := a.periodRepository.FindWithJoin(
		map[string]any{
			"id": UUIDPeriod,
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	timeLoc, _ := time.LoadLocation(config.GetEnv("TIMEZONE", "Asia/Jakarta"))

	if len(findAttendancePeriod) == 0 {
		return nil, errors.New("period not found")
	}

	output.Period.StartDate = findAttendancePeriod[0].StartDate.In(timeLoc)
	output.Period.EndDate = findAttendancePeriod[0].EndDate.In(timeLoc)

	totalWeekday := 0
	totalDays := 0
	for date := findAttendancePeriod[0].StartDate; !date.After(findAttendancePeriod[0].EndDate.AddDate(0, 0, 1)); date = date.AddDate(0, 0, 1) {
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			totalWeekday++
		}

		totalDays++
	}

	// count attendance
	arrDetailAttend := make([]time.Time, 0)
	findAttendance, _ := a.attendanceRepository.FindWithJoin(
		map[string]any{
			"user_id":   UUIDUser,
			"period_id": UUIDPeriod,
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	sort.Slice(findAttendance, func(i, j int) bool {
		return findAttendance[i].DateIn.Before(findAttendance[j].DateIn)
	})

	for _, atd := range findAttendance {
		arrDetailAttend = append(arrDetailAttend, atd.DateIn.In(timeLoc))
	}

	output.Attendances.TotalShouldAttend = totalDays
	output.Attendances.SalaryDay = float64(findUser.Salary / totalDays)
	output.Attendances.TotalAttend = len(findAttendance)
	output.Attendances.Formula = "Salary Day x Total Attend"
	output.Attendances.DetailAttend = arrDetailAttend

	// find overtime
	salaryHourly := float64(findUser.Salary/totalDays) / 8
	totalOvertime := 0
	arrDetailOvertime := make([]shared.DetailOvertimes, 0)
	findOvertime, _ := a.overtimeRepository.FindWithJoin(
		map[string]any{
			"user_id":   UUIDUser,
			"period_id": UUIDPeriod,
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	sort.Slice(findOvertime, func(i, j int) bool {
		return findOvertime[i].Date.Before(findOvertime[j].Date)
	})

	for _, ovt := range findOvertime {
		arrDetailOvertime = append(arrDetailOvertime, shared.DetailOvertimes{
			Day:   ovt.Date.In(timeLoc),
			Hours: ovt.Hours,
		})
		totalOvertime += ovt.Hours
	}

	output.Overtimes.SalaryHours = salaryHourly
	output.Overtimes.Formula = "Salary Hourly x Total Overtime x 2"
	output.Overtimes.OvertimeSalary = salaryHourly * float64(totalOvertime) * 2
	output.Overtimes.TotalOvertime = totalOvertime
	output.Overtimes.DetailOvertimes = arrDetailOvertime

	// find reimburse
	totalReimburse := float64(0)
	arrReimburse := make([]shared.Reimbursement, 0)
	findReimburse, _ := a.reimbursementRepository.FindWithJoin(
		map[string]any{
			"user_id":   UUIDUser,
			"period_id": UUIDPeriod,
		},
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	sort.Slice(findReimburse, func(i, j int) bool {
		return findReimburse[i].CreatedAt.Before(findReimburse[j].CreatedAt)
	})

	for _, ovt := range findReimburse {
		arrReimburse = append(arrReimburse, shared.Reimbursement{
			Amount:      ovt.Amount,
			Description: ovt.Description,
		})
		totalReimburse += ovt.Amount
	}
	output.Reimbursements = arrReimburse
	output.TakeHomePay = (output.Attendances.SalaryDay * float64(
		output.Attendances.TotalAttend,
	)) + output.Overtimes.OvertimeSalary + totalReimburse

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		UserID:    UUIDUser,
		Action:    "get_payslip",
		Entity:    "PaySlip",
		EntityID:  UUIDUser,
		IpAddress: *ipAddress,
		RequestID: *requestID,
		PeriodID:  &UUIDPeriod,
	})

	_, err = a.logAuditRepository.Create(ctx, arrLogAudit)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return output, nil
}

func NewPayslipService(
	payslipRepository repositories.PayslipRepositoryInterface,
	logAuditRepository repositories.LogAuditRepositoryInterface,
	payrollRepository repositories.PayrollRepositoryInterface,
	userRepository repositories.UserRepositoryInterface,
	periodRepository repositories.AttendancePeriodRepositoryInterface,
	attendanceRepository repositories.AttendanceRepositoryInterface,
	overtimeRepository repositories.OvertimeRepositoryInterface,
	reimbursementRepository repositories.ReimbursementRepositoryInterface,
) PayslipServiceInterface {
	return &payslipService{
		payslipRepository:       payslipRepository,
		payrollRepository:       payrollRepository,
		logAuditRepository:      logAuditRepository,
		userRepository:          userRepository,
		periodRepository:        periodRepository,
		attendanceRepository:    attendanceRepository,
		overtimeRepository:      overtimeRepository,
		reimbursementRepository: reimbursementRepository,
	}
}
