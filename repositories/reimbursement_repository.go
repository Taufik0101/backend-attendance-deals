package repositories

import (
	"backend-attendance-deals/database/models"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReimbursementRepositoryInterface interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Create(ctx context.Context, reimbursements []*models.Reimbursement) ([]*models.Reimbursement, error)
	FindWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) ([]*models.Reimbursement, error)
	FindOneWithJoin(
		whereClause map[string]any, // key is where column, value is the value to filter
		whereNotClause map[string]any, // key is where column, value is the value to filter
		orders []string,
		joins []string,
		groups []string,
		relations []string,
	) (*models.Reimbursement, error)
}

type reimbursementRepository struct {
	db *gorm.DB
}

func (a reimbursementRepository) FindOneWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) (*models.Reimbursement, error) {
	var reimbursement *models.Reimbursement

	query := a.db.Model(&models.Reimbursement{})

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

	err := query.First(&reimbursement).Error
	if err != nil {
		log.Errorf("[ReimbursementRepository][Find]: %v", err)
		return nil, err
	}

	return reimbursement, err
}

func (a reimbursementRepository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (a reimbursementRepository) Create(ctx context.Context, reimbursements []*models.Reimbursement) ([]*models.Reimbursement, error) {
	if len(reimbursements) < 1 {
		return make([]*models.Reimbursement, 0), nil
	}

	db := ExtractTx(ctx, a.db)

	err := db.Save(&reimbursements).Error
	if err != nil {
		log.Errorf("[ReimbursementRepository][Create] failed to exec Save, error: %v", err)
	}
	return reimbursements, err
}

func (a reimbursementRepository) FindWithJoin(whereClause map[string]any, whereNotClause map[string]any, orders []string, joins []string, groups []string, relations []string) ([]*models.Reimbursement, error) {
	var reimbursements []*models.Reimbursement

	query := a.db.Model(&models.Reimbursement{})

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

	err := query.Find(&reimbursements).Error
	if err != nil {
		log.Errorf("[ReimbursementRepository][Find]: %v", err)
		return nil, err
	}

	return reimbursements, err
}

func NewReimbursementRepository(db *gorm.DB) ReimbursementRepositoryInterface {
	return &reimbursementRepository{db: db}
}
