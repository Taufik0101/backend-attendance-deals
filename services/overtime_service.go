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
	"github.com/sirupsen/logrus"
	"strings"
	"time"

	"github.com/jinzhu/now"
)

type OvertimeServiceInterface interface {
	Create(ctx context.Context, input shared.CreateOvertimeInput) (*models.Overtime, error)
}

type overtimeService struct {
	overtimeRepository   repositories.OvertimeRepositoryInterface
	attendanceRepository repositories.AttendanceRepositoryInterface
	logAuditRepository   repositories.LogAuditRepositoryInterface
}

func (a overtimeService) Create(ctx context.Context, input shared.CreateOvertimeInput) (*models.Overtime, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)
	currentUser := utils.GetCurrentUser(&ctx)
	UUIDUser, _ := uuid.Parse(currentUser.UserId)

	if strings.ToLower(currentUser.UserType) != "employee" {
		return nil, errors.New("FORBIDDEN")
	}

	output := new(models.Overtime)
	timeLoc, _ := time.LoadLocation(config.GetEnv("TIMEZONE", "Asia/Jakarta"))
	tNow := now.With(time.Now()).In(timeLoc)

	nowBeginningDay := now.With(tNow).BeginningOfDay()
	nowEndDay := now.With(tNow).EndOfDay()
	// check already check
	checkAlreadyCheckOut, _ := a.attendanceRepository.FindWithJoin(
		map[string]any{
			"date_out >= ?": nowBeginningDay,
			"date_out < ?":  nowEndDay,
			"user_id":       UUIDUser,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(checkAlreadyCheckOut) == 0 {
		return nil, errors.New("you must checkout first")
	}

	if input.Hours > 3 {
		return nil, errors.New("max overtime is 3 hours each day")
	}

	// already submit overtime
	checkAlreadyOvertime, _ := a.overtimeRepository.FindWithJoin(
		map[string]any{
			"date >= ?": nowBeginningDay,
			"date <= ?": nowEndDay,
			"user_id":   UUIDUser,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(checkAlreadyOvertime) > 0 {
		return nil, errors.New("only submit 1 overtime for each day")
	}

	arrCreateOvertime := make([]*models.Overtime, 0)
	arrCreateOvertime = append(arrCreateOvertime, &models.Overtime{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		Date:      tNow,
		Hours:     input.Hours,
		UserID:    UUIDUser,
		IpAddress: *ipAddress,
		PeriodID:  checkAlreadyCheckOut[0].PeriodID,
	})

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		UserID:    UUIDUser,
		Action:    "submit_overtime",
		Entity:    "Overtime",
		EntityID:  UUIDUser,
		IpAddress: *ipAddress,
		RequestID: *requestID,
		PeriodID:  &checkAlreadyCheckOut[0].PeriodID,
	})

	errTx := a.overtimeRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		createOvertime, err := a.overtimeRepository.Create(ctx, arrCreateOvertime)

		if err != nil {
			return errors.New(err.Error())
		}

		_, err = a.logAuditRepository.Create(ctx, arrLogAudit)

		if err != nil {
			return errors.New(err.Error())
		}

		output = createOvertime[0]

		return nil
	})

	if errTx != nil {
		logrus.Errorf("[overtimeService][Create] failed to commit transaction, error: %v", errTx)
		return nil, errTx
	}

	return output, nil
}

func NewOvertimeService(
	overtimeRepository repositories.OvertimeRepositoryInterface,
	attendanceRepository repositories.AttendanceRepositoryInterface,
	logAuditRepository repositories.LogAuditRepositoryInterface,
) OvertimeServiceInterface {
	return &overtimeService{
		overtimeRepository:   overtimeRepository,
		attendanceRepository: attendanceRepository,
		logAuditRepository:   logAuditRepository,
	}
}
