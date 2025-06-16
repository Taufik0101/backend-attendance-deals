package repositories

import (
	"backend-attendance-deals/database/models"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PayslipRepositoryInterface interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Create(ctx context.Context, payslips []*models.Payslip) ([]*models.Payslip, error)
	FindWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) ([]*models.Payslip, error)
	FindOneWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) (*models.Payslip, error)
}

type payslipRepository struct {
	db *gorm.DB
}

func (a payslipRepository) FindOneWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) (*models.Payslip, error) {
	var payslip *models.Payslip

	query := a.db.Model(&models.Payslip{})

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

	err := query.First(&payslip).Error
	if err != nil {
		log.Errorf("[PayslipRepository][Find]: %v", err)
		return nil, err
	}

	return payslip, err
}

func (a payslipRepository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (a payslipRepository) Create(ctx context.Context, payslips []*models.Payslip) ([]*models.Payslip, error) {
	if len(payslips) < 1 {
		return make([]*models.Payslip, 0), nil
	}

	db := ExtractTx(ctx, a.db)

	err := db.Save(&payslips).Error
	if err != nil {
		log.Errorf("[PayslipRepository][Create] failed to exec Save, error: %v", err)
	}
	return payslips, err
}

func (a payslipRepository) FindWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) ([]*models.Payslip, error) {
	var payslips []*models.Payslip

	query := a.db.Model(&models.Payslip{})

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

	err := query.Find(&payslips).Error
	if err != nil {
		log.Errorf("[PayslipRepository][Find]: %v", err)
		return nil, err
	}

	return payslips, err
}

func NewPayslipRepository(db *gorm.DB) PayslipRepositoryInterface {
	return &payslipRepository{db: db}
}
