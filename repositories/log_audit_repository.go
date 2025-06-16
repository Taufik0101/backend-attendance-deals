package repositories

import (
	"backend-attendance-deals/database/models"
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type LogAuditRepositoryInterface interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
	Create(ctx context.Context, logAudits []*models.LogAudit) ([]*models.LogAudit, error)
}

type logAuditRepository struct {
	db *gorm.DB
}

func (l logAuditRepository) Create(ctx context.Context, logAudits []*models.LogAudit) ([]*models.LogAudit, error) {
	if len(logAudits) < 1 {
		return make([]*models.LogAudit, 0), nil
	}

	db := ExtractTx(ctx, l.db)

	err := db.Save(&logAudits).Error
	if err != nil {
		log.Errorf("[LogAuditRepository][Create] failed to exec Save, error: %v", err)
	}
	return logAudits, err
}

func (l logAuditRepository) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx := l.db.WithContext(ctx).Begin()

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

func NewLogAuditRepository(db *gorm.DB) LogAuditRepositoryInterface {
	return &logAuditRepository{db: db}
}
