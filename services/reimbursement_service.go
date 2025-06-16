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
	"github.com/jinzhu/now"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ReimbursementServiceInterface interface {
	Create(ctx context.Context, input shared.CreateReimbursementInput) (*models.Reimbursement, error)
}

type reimbursementService struct {
	reimbursementRepository    repositories.ReimbursementRepositoryInterface
	logAuditRepository         repositories.LogAuditRepositoryInterface
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface
}

func (a reimbursementService) Create(ctx context.Context, input shared.CreateReimbursementInput) (*models.Reimbursement, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)
	currentUser := utils.GetCurrentUser(&ctx)
	UUIDUser, _ := uuid.Parse(currentUser.UserId)

	if strings.ToLower(currentUser.UserType) != "employee" {
		return nil, errors.New("FORBIDDEN")
	}

	output := new(models.Reimbursement)
	timeLoc, _ := time.LoadLocation(config.GetEnv("TIMEZONE", "Asia/Jakarta"))
	Date, err := time.Parse("2006-01-02", input.DayActivity)

	if err != nil {
		return nil, err
	}

	Date = now.With(Date).In(timeLoc)
	tNowDay := Date.Weekday()

	if tNowDay == time.Saturday || tNowDay == time.Sunday {
		return nil, errors.New("cannot submit reimburse because its not on weekday")
	}

	// check start date inside at existing period
	checkPeriodActive, _ := a.attendancePeriodRepository.FindWithJoin(
		map[string]any{
			"start_date <= ?": Date,
			"end_date >= ?":   Date,
		},
		nil,
		nil,
		nil,
		nil,
		nil)

	if len(checkPeriodActive) == 0 {
		return nil, errors.New("day activity is not on period")
	}

	arrCreateReimbursement := make([]*models.Reimbursement, 0)
	arrCreateReimbursement = append(arrCreateReimbursement, &models.Reimbursement{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		Amount:      input.Amount,
		Description: input.Description,
		UserID:      UUIDUser,
		IpAddress:   *ipAddress,
		PeriodID:    checkPeriodActive[0].ID,
	})

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &UUIDUser,
		},
		UserID:    UUIDUser,
		Action:    "submit_reimbursement",
		Entity:    "Reimbursement",
		EntityID:  UUIDUser,
		IpAddress: *ipAddress,
		RequestID: *requestID,
		PeriodID:  &checkPeriodActive[0].ID,
	})

	errTx := a.reimbursementRepository.WithinTransaction(ctx, func(ctx context.Context) error {

		createReimbursement, err := a.reimbursementRepository.Create(ctx, arrCreateReimbursement)

		if err != nil {
			return errors.New(err.Error())
		}

		_, err = a.logAuditRepository.Create(ctx, arrLogAudit)

		if err != nil {
			return errors.New(err.Error())
		}

		output = createReimbursement[0]

		return nil
	})

	if errTx != nil {
		logrus.Errorf("[reimbursementService][Create] failed to commit transaction, error: %v", errTx)
		return nil, errTx
	}

	return output, nil
}

func NewReimbursementService(
	reimbursementRepository repositories.ReimbursementRepositoryInterface,
	attendancePeriodRepository repositories.AttendancePeriodRepositoryInterface,
	logAuditRepository repositories.LogAuditRepositoryInterface,
) ReimbursementServiceInterface {
	return &reimbursementService{
		reimbursementRepository:    reimbursementRepository,
		attendancePeriodRepository: attendancePeriodRepository,
		logAuditRepository:         logAuditRepository,
	}
}
