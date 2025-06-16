package repositories

import (
	"backend-attendance-deals/database/models"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendancePeriodRepositoryInterface interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Create(ctx context.Context, attendancePeriods []*models.AttendancePeriod) ([]*models.AttendancePeriod, error)
	FindWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) ([]*models.AttendancePeriod, error)
}

type attendancePeriodRepository struct {
	db *gorm.DB
}

func (a attendancePeriodRepository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx := a.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		log.Errorf("failed to begin transaction: %v", tx.Error)
		return tx.Error
	}

	if err := tFunc(InjectTx(ctx, tx)); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Errorf("failed to commit transaction: %v", tx.Error)
		return err
	}

	return nil
}

func (a attendancePeriodRepository) Create(ctx context.Context, attendancePeriods []*models.AttendancePeriod) ([]*models.AttendancePeriod, error) {
	if len(attendancePeriods) < 1 {
		return make([]*models.AttendancePeriod, 0), nil
	}

	db := ExtractTx(ctx, a.db)

	err := db.Save(&attendancePeriods).Error
	if err != nil {
		log.Errorf("[AttendancePeriodRepository][Create] failed to exec Save, error: %v", err)
	}
	return attendancePeriods, err
}

func (a attendancePeriodRepository) FindWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) ([]*models.AttendancePeriod, error) {
	var attendancePeriods []*models.AttendancePeriod

	query := a.db.Model(&models.AttendancePeriod{})

	for _, join := range joins {
		query.Joins(join)
	}

	for key, value := range whereClause {
		query = query.Where(key, value)
	}

	for key, value := range whereNotClause {
		query = query.Not(key, value)
	}

	for _, orderClause := range orders {
		query = query.Order(orderClause)
	}

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	for _, groupClause := range groups {
		query.Group(groupClause)
	}

	err := query.Find(&attendancePeriods).Error
	if err != nil {
		log.Errorf("[AttendancePeriodRepository][Find]: %v", err)
		return nil, err
	}

	return attendancePeriods, err
}

func NewAttendancePeriodRepository(db *gorm.DB) AttendancePeriodRepositoryInterface {
	return &attendancePeriodRepository{db: db}
}
