package repositories

import (
	"backend-attendance-deals/database/models"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AttendanceRepositoryInterface interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Create(ctx context.Context, attendances []*models.Attendance) ([]*models.Attendance, error)
	FindWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) ([]*models.Attendance, error)
	FindOneWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) (*models.Attendance, error)
}

type attendanceRepository struct {
	db *gorm.DB
}

func (a attendanceRepository) FindOneWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) (*models.Attendance, error) {
	var attendance *models.Attendance

	query := a.db.Model(&models.Attendance{})

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

	err := query.First(&attendance).Error
	if err != nil {
		log.Errorf("[AttendanceRepository][Find]: %v", err)
		return nil, err
	}

	return attendance, err
}

func (a attendanceRepository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (a attendanceRepository) Create(ctx context.Context, attendances []*models.Attendance) ([]*models.Attendance, error) {
	if len(attendances) < 1 {
		return make([]*models.Attendance, 0), nil
	}

	db := ExtractTx(ctx, a.db)

	err := db.Save(&attendances).Error
	if err != nil {
		log.Errorf("[AttendanceRepository][Create] failed to exec Save, error: %v", err)
	}
	return attendances, err
}

func (a attendanceRepository) FindWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) ([]*models.Attendance, error) {
	var attendances []*models.Attendance

	query := a.db.Model(&models.Attendance{})

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

	err := query.Find(&attendances).Error
	if err != nil {
		log.Errorf("[AttendanceRepository][Find]: %v", err)
		return nil, err
	}

	return attendances, err
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepositoryInterface {
	return &attendanceRepository{db: db}
}
