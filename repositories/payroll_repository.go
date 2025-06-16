package repositories

import (
	"backend-attendance-deals/database/models"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PayrollRepositoryInterface interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Create(ctx context.Context, payrolls []*models.Payroll) ([]*models.Payroll, error)
	FindWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) ([]*models.Payroll, error)
	FindOneWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) (*models.Payroll, error)
}

type payrollRepository struct {
	db *gorm.DB
}

func (a payrollRepository) FindOneWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) (*models.Payroll, error) {
	var payroll *models.Payroll

	query := a.db.Model(&models.Payroll{})

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

	err := query.First(&payroll).Error
	if err != nil {
		log.Errorf("[PayrollRepository][Find]: %v", err)
		return nil, err
	}

	return payroll, err
}

func (a payrollRepository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (a payrollRepository) Create(ctx context.Context, payrolls []*models.Payroll) ([]*models.Payroll, error) {
	if len(payrolls) < 1 {
		return make([]*models.Payroll, 0), nil
	}

	db := ExtractTx(ctx, a.db)

	err := db.Save(&payrolls).Error
	if err != nil {
		log.Errorf("[PayrollRepository][Create] failed to exec Save, error: %v", err)
	}
	return payrolls, err
}

func (a payrollRepository) FindWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) ([]*models.Payroll, error) {
	var payrolls []*models.Payroll

	query := a.db.Model(&models.Payroll{})

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

	err := query.Find(&payrolls).Error
	if err != nil {
		log.Errorf("[PayrollRepository][Find]: %v", err)
		return nil, err
	}

	return payrolls, err
}

func NewPayrollRepository(db *gorm.DB) PayrollRepositoryInterface {
	return &payrollRepository{db: db}
}
