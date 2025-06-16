package repositories

import (
	"backend-attendance-deals/database/models"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OvertimeRepositoryInterface interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Create(ctx context.Context, overtimes []*models.Overtime) ([]*models.Overtime, error)
	FindWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) ([]*models.Overtime, error)
	FindOneWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) (*models.Overtime, error)
}

type overtimeRepository struct {
	db *gorm.DB
}

func (a overtimeRepository) FindOneWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) (*models.Overtime, error) {
	var overtime *models.Overtime

	query := a.db.Model(&models.Overtime{})

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

	err := query.First(&overtime).Error
	if err != nil {
		log.Errorf("[OvertimeRepository][Find]: %v", err)
		return nil, err
	}

	return overtime, err
}

func (a overtimeRepository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (a overtimeRepository) Create(ctx context.Context, overtimes []*models.Overtime) ([]*models.Overtime, error) {
	if len(overtimes) < 1 {
		return make([]*models.Overtime, 0), nil
	}

	db := ExtractTx(ctx, a.db)

	err := db.Save(&overtimes).Error
	if err != nil {
		log.Errorf("[OvertimeRepository][Create] failed to exec Save, error: %v", err)
	}
	return overtimes, err
}

func (a overtimeRepository) FindWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) ([]*models.Overtime, error) {
	var overtimes []*models.Overtime

	query := a.db.Model(&models.Overtime{})

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

	err := query.Find(&overtimes).Error
	if err != nil {
		log.Errorf("[OvertimeRepository][Find]: %v", err)
		return nil, err
	}

	return overtimes, err
}

func NewOvertimeRepository(db *gorm.DB) OvertimeRepositoryInterface {
	return &overtimeRepository{db: db}
}
