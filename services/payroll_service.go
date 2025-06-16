package services

import (
	"backend-attendance-deals/database/models"
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/pkg/utils"
	"backend-attendance-deals/repositories"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type PayrollServiceInterface interface {
	Create(ctx context.Context, input shared.CreatePayrollInput) (*models.Payroll, error)
}

type payrollService struct {
	payrollRepository          repositories.PayrollRepositoryInterface
	attendanceRepository       repositories.AttendanceRepositoryInterface
	overtimeRepository         repositories.OvertimeRepositoryInterface
	reimbursementRepository    repositories.ReimbursementRepositoryInterface
	logAuditRepository         repositories.LogAuditRepositoryInterface
	userRepository             repositories.UserRepositoryInterface
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface
	payslipRepository          repositories.PayslipRepositoryInterface
}

func (a payrollService) Create(ctx context.Context, input shared.CreatePayrollInput) (*models.Payroll, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)
	currentUser := utils.GetCurrentUser(&ctx)
	UUIDUser, _ := uuid.Parse(currentUser.UserId)
	UUIDPeriod, _ := uuid.Parse(input.PeriodID)

	if strings.ToLower(currentUser.UserType) != "admin" {
		return nil, errors.New("FORBIDDEN")
	}

	output := new(models.Payroll)

	// check period id
	findAttendancePeriod, _ := a.attendancePeriodRepository.FindWithJoin(
		map[string]any{
			"id": input.PeriodID,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(findAttendancePeriod) == 0 {
		return nil, errors.New("attendance period not found")
	}

	totalWeekday := 0
	totalDays := 0
	for date := findAttendancePeriod[0].StartDate; !date.After(findAttendancePeriod[0].EndDate.AddDate(0, 0, 1)); date = date.AddDate(0, 0, 1) {
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			totalWeekday++
		}

		totalDays++
	}

	// check payroll ald run
	findPayroll, _ := a.payrollRepository.FindOneWithJoin(
		map[string]any{
			"period_id": input.PeriodID,
		},
		nil, nil, nil, nil, nil)

	if findPayroll != nil {
		return nil, errors.New("payroll only can run once")
	}

	// get all attendance
	findAttendance, _ := a.attendanceRepository.FindWithJoin(
		map[string]any{
			"period_id": input.PeriodID,
		},
		nil, nil, nil, nil, nil)

	// find overtime
	findOvertime, _ := a.overtimeRepository.FindWithJoin(
		map[string]any{
			"period_id": input.PeriodID,
		},
		nil, nil, nil, nil, nil)

	// find reimbursement
	findReimbursement, _ := a.reimbursementRepository.FindWithJoin(
		map[string]any{
			"period_id": input.PeriodID,
		},
		nil, nil, nil, nil, nil)

	// find user employee
	findUser, _ := a.userRepository.FindWithAttribute(map[string]any{
		"role": "employee",
	}, nil, nil)

	// calculate
	mapTotalAttendance := make(map[uuid.UUID]int)
	mapTotalOvertime := make(map[uuid.UUID]int)
	mapTotalReimbursement := make(map[uuid.UUID]float64)

	for _, atd := range findAttendance {
		if _, exists := mapTotalAttendance[atd.UserID]; !exists {
			temp := mapTotalAttendance[atd.UserID]
			temp = 1
			mapTotalAttendance[atd.UserID] = temp
		} else {
			temp := mapTotalAttendance[atd.UserID]
			temp++
			mapTotalAttendance[atd.UserID] = temp
		}
	}

	for _, atd := range findOvertime {
		if _, exists := mapTotalOvertime[atd.UserID]; !exists {
			temp := mapTotalOvertime[atd.UserID]
			temp = atd.Hours
			mapTotalOvertime[atd.UserID] = temp
		} else {
			temp := mapTotalOvertime[atd.UserID]
			temp += atd.Hours
			mapTotalOvertime[atd.UserID] = temp
		}
	}

	for _, atd := range findReimbursement {
		if _, exists := mapTotalReimbursement[atd.UserID]; !exists {
			temp := mapTotalReimbursement[atd.UserID]
			temp = atd.Amount
			mapTotalReimbursement[atd.UserID] = temp
		} else {
			temp := mapTotalReimbursement[atd.UserID]
			temp += atd.Amount
			mapTotalReimbursement[atd.UserID] = temp
		}
	}

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		UserID:    UUIDUser,
		Action:    "submit_payroll",
		Entity:    "Payroll",
		EntityID:  UUIDUser,
		IpAddress: *ipAddress,
		RequestID: *requestID,
		PeriodID:  &UUIDPeriod,
	})

	errTx := a.payrollRepository.WithinTransaction(ctx, func(ctx context.Context) error {
		arrCreatePayroll := make([]*models.Payroll, 0)
		arrCreatePayroll = append(arrCreatePayroll, &models.Payroll{
			BaseModel: models.BaseModel{
				CreatedBy: &UUIDUser,
			},
			PeriodID: UUIDPeriod,
		})

		createPayroll, err := a.payrollRepository.Create(ctx, arrCreatePayroll)

		if err != nil {
			return errors.New(err.Error())
		}

		// make payslip
		arrPaySlip := make([]*models.Payslip, 0)
		for _, user := range findUser {
			salaryDaily := user.Salary / totalDays
			salaryHourly := salaryDaily / 8

			salaryAttendance := salaryDaily * mapTotalAttendance[user.ID]
			salaryOvertime := (salaryHourly * mapTotalOvertime[user.ID]) * 2
			totalPay := float64(salaryAttendance) + float64(salaryOvertime) + mapTotalReimbursement[user.ID]

			arrPaySlip = append(arrPaySlip, &models.Payslip{
				BaseModel: models.BaseModel{
					CreatedBy: &UUIDUser,
				},
				AttendanceDays: mapTotalAttendance[user.ID],
				OvertimeHours:  mapTotalOvertime[user.ID],
				Reimbursement:  mapTotalReimbursement[user.ID],
				BaseSalary:     float64(user.Salary),
				TotalPay:       totalPay,
				UserID:         user.ID,
				PayrollID:      createPayroll[0].ID,
			})
		}

		_, err = a.payslipRepository.Create(ctx, arrPaySlip)

		if err != nil {
			return errors.New(err.Error())
		}

		_, err = a.logAuditRepository.Create(ctx, arrLogAudit)

		if err != nil {
			return errors.New(err.Error())
		}

		output = createPayroll[0]

		return nil
	})

	if errTx != nil {
		logrus.Errorf("[payrollService][Create] failed to commit transaction, error: %v", errTx)
		return nil, errTx
	}

	return output, nil
}

func NewPayrollService(
	payrollRepository repositories.PayrollRepositoryInterface,
	attendanceRepository repositories.AttendanceRepositoryInterface,
	overtimeRepository repositories.OvertimeRepositoryInterface,
	reimbursementRepository repositories.ReimbursementRepositoryInterface,
	logAuditRepository repositories.LogAuditRepositoryInterface,
	userRepository repositories.UserRepositoryInterface,
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface,
	payslipRepository repositories.PayslipRepositoryInterface,
) PayrollServiceInterface {
	return &payrollService{
		payrollRepository:          payrollRepository,
		attendanceRepository:       attendanceRepository,
		overtimeRepository:         overtimeRepository,
		reimbursementRepository:    reimbursementRepository,
		logAuditRepository:         logAuditRepository,
		userRepository:             userRepository,
		attendancePeriodRepository: attendancePeriodRepository,
		payslipRepository:          payslipRepository,
	}
}
