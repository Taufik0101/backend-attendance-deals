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

type AttendancePeriodServiceInterface interface {
	Create(ctx context.Context, input shared.CreateAttendancePeriodInput) (*models.AttendancePeriod, error)
}

type attendancePeriodService struct {
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface
	logAuditRepository         repositories.LogAuditRepositoryInterface
}

func (a attendancePeriodService) Create(ctx context.Context, input shared.CreateAttendancePeriodInput) (*models.AttendancePeriod, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)
	currentUser := utils.GetCurrentUser(&ctx)
	UUIDUser, _ := uuid.Parse(currentUser.UserId)

	if strings.ToLower(currentUser.UserType) != "admin" {
		return nil, errors.New("FORBIDDEN")
	}

	output := new(models.AttendancePeriod)
	timeLoc, _ := time.LoadLocation(config.GetEnv("TIMEZONE", "Asia/Jakarta"))
	startDate, err := time.Parse("2006-01-02", input.StartDate)

	if err != nil {
		return nil, err
	}

	// check if first date of month or not
	if startDate.Day() != 1 {
		return nil, errors.New("start date of month must be 1")
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)

	if err != nil {
		return nil, err
	}

	// check if end date of month or not
	if endDate.AddDate(0, 0, 1).Day() != 1 {
		return nil, errors.New("end date is not end of date")
	}

	startDate = startDate.In(timeLoc)
	endDate = endDate.In(timeLoc)

	startDate = now.With(startDate).BeginningOfDay()
	endDate = now.With(endDate).EndOfDay()

	// check start date inside at existing period
	checkStartExistingPeriod, _ := a.attendancePeriodRepository.FindWithJoin(
		map[string]any{
			"start_date <= ?": startDate,
			"end_date >= ?":   startDate,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(checkStartExistingPeriod) > 0 {
		return nil, errors.New("start date inside period active")
	}

	// check end date inside at existing period
	checkEndExistingPeriod, _ := a.attendancePeriodRepository.FindWithJoin(
		map[string]any{
			"start_date <= ?": endDate,
			"end_date >= ?":   endDate,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(checkEndExistingPeriod) > 0 {
		return nil, errors.New("end date inside period active")
	}

	// check start end date inside at existing period
	checkStartEndExistingPeriod, _ := a.attendancePeriodRepository.FindWithJoin(
		map[string]any{
			"start_date >= ?": startDate,
			"end_date <= ?":   endDate,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(checkStartEndExistingPeriod) > 0 {
		return nil, errors.New("start end date inside period active")
	}

	arrAttendancePeriod := make([]*models.AttendancePeriod, 0)
	arrAttendancePeriod = append(arrAttendancePeriod, &models.AttendancePeriod{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		StartDate: startDate,
		EndDate:   endDate,
	})

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		UserID:    UUIDUser,
		Action:    "create_attendance_period",
		Entity:    "Attendance Period",
		EntityID:  UUIDUser,
		IpAddress: *ipAddress,
		RequestID: *requestID,
	})

	errTx := a.attendancePeriodRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		createAttendancePeriod, err := a.attendancePeriodRepository.Create(ctx, arrAttendancePeriod)

		if err != nil {
			return errors.New(err.Error())
		}

		_, err = a.logAuditRepository.Create(ctx, arrLogAudit)

		if err != nil {
			return errors.New(err.Error())
		}

		output = createAttendancePeriod[0]

		return nil
	})

	if errTx != nil {
		logrus.Errorf("[attendancePeriodService][Create] failed to commit transaction, error: %v", errTx)
		return nil, errTx
	}

	return output, nil
}

func NewAttendancePeriodService(
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface,
	logAuditRepository repositories.LogAuditRepositoryInterface,
) AttendancePeriodServiceInterface {
	return &attendancePeriodService{
		attendancePeriodRepository: attendancePeriodRepository,
		logAuditRepository:         logAuditRepository,
	}
}
