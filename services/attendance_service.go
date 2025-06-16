package services

import (
	"backend-attendance-deals/config"
	"backend-attendance-deals/database/models"
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/pkg/utils"
	"backend-attendance-deals/repositories"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"slices"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

type AttendanceServiceInterface interface {
	Create(ctx context.Context, input shared.CreateAttendanceInput) (*models.Attendance, error)
}

type attendanceService struct {
	attendanceRepository       repositories.AttendanceRepositoryInterface
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface
	logAuditRepository         repositories.LogAuditRepositoryInterface
	payrollRepository          repositories.PayrollRepositoryInterface
}

func (a attendanceService) Create(ctx context.Context, input shared.CreateAttendanceInput) (*models.Attendance, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)
	currentUser := utils.GetCurrentUser(&ctx)
	UUIDUser, _ := uuid.Parse(currentUser.UserId)

	if strings.ToLower(currentUser.UserType) != "employee" {
		return nil, errors.New("FORBIDDEN")
	}

	if !slices.Contains([]string{"in", "out"}, strings.ToLower(input.Type)) {
		return nil, errors.New("only type check in or check out")
	}

	output := new(models.Attendance)
	timeLoc, _ := time.LoadLocation(config.GetEnv("TIMEZONE", "Asia/Jakarta"))
	tNow := now.With(time.Now()).In(timeLoc)
	tNowDay := tNow.Weekday()

	if tNowDay == time.Saturday || tNowDay == time.Sunday {
		return nil, errors.New("cannot attendance on weekday")
	}

	// check start date inside at existing period
	checkPeriodActive, _ := a.attendancePeriodRepository.FindWithJoin(
		map[string]any{
			"start_date <= ?": tNow,
			"end_date >= ?":   tNow,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(checkPeriodActive) == 0 {
		return nil, errors.New("no period active")
	}

	// check attendance period already processed or not
	checkAttendanceProcessed, _ := a.payrollRepository.FindOneWithJoin(
		map[string]any{
			"period_id": checkPeriodActive[0],
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if checkAttendanceProcessed != nil {
		return nil, errors.New("attendance period already processed")
	}

	nowBeginningDay := now.With(tNow).BeginningOfDay()
	nowEndDay := now.With(tNow).EndOfDay()
	// check already check
	checkAlreadyCheck, _ := a.attendanceRepository.FindWithJoin(
		map[string]any{
			"date_in >= ?": nowBeginningDay,
			"date_in < ?":  nowEndDay,
			"user_id":      UUIDUser,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	arrCreateAttendance := make([]*models.Attendance, 0)
	if len(checkAlreadyCheck) == 0 {

		if strings.ToLower(input.Type) == "out" {
			return nil, errors.New("you must check in first")
		}

		arrCreateAttendance = append(arrCreateAttendance, &models.Attendance{
			BaseModel: models.BaseModel{
				CreatedBy: &UUIDUser,
			},
			DateIn:    tNow,
			UserID:    UUIDUser,
			PeriodID:  checkPeriodActive[0].ID,
			IpAddress: *ipAddress,
		})
	} else {
		if strings.ToLower(input.Type) == "out" && checkAlreadyCheck[0].DateOut == nil {
			checkAlreadyCheck[0].BaseModel.UpdatedBy = &UUIDUser
			checkAlreadyCheck[0].DateOut = &tNow
			arrCreateAttendance = append(arrCreateAttendance, checkAlreadyCheck[0])
		}
	}

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		UserID:    UUIDUser,
		Action:    fmt.Sprintf("submit_attendance_%s", strings.ToLower(input.Type)),
		Entity:    "Attendance",
		EntityID:  UUIDUser,
		IpAddress: *ipAddress,
		RequestID: *requestID,
		PeriodID:  &checkPeriodActive[0].ID,
	})

	errTx := a.attendanceRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		if len(arrCreateAttendance) > 0 {
			check, err := a.attendanceRepository.Create(ctx, arrCreateAttendance)

			if err != nil {
				return errors.New(err.Error())
			}

			checkAlreadyCheck = check
		}

		_, err := a.logAuditRepository.Create(ctx, arrLogAudit)

		if err != nil {
			return errors.New(err.Error())
		}

		output = checkAlreadyCheck[0]

		return nil
	})

	if errTx != nil {
		logrus.Errorf("[attendanceService][Create] failed to commit transaction, error: %v", errTx)
		return nil, errTx
	}

	return output, nil
}

func NewAttendanceService(
	attendanceRepository repositories.AttendanceRepositoryInterface,
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface,
	logAuditRepository repositories.LogAuditRepositoryInterface,
	payrollRepository repositories.PayrollRepositoryInterface,
) AttendanceServiceInterface {
	return &attendanceService{
		attendanceRepository:       attendanceRepository,
		attendancePeriodRepository: attendancePeriodRepository,
		logAuditRepository:         logAuditRepository,
		payrollRepository:          payrollRepository,
	}
}
