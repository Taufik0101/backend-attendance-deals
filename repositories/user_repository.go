package repositories

import (
	"backend-attendance-deals/database/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	FindOneWithAttribute(whereClause map[string]any, whereNotClause map[string]any, relations []string) (*models.User, error)
	FindWithAttribute(whereClause map[string]any, whereNotClause map[string]any, relations []string) ([]*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) FindWithAttribute(whereClause map[string]any, whereNotClause map[string]any, relations []string) ([]*models.User, error) {
	var user []*models.User

	query := u.db

	for whereQuery, whereArgs := range whereClause {
		query = query.Where(whereQuery, whereArgs)
	}

	for whereQuery, whereArgs := range whereNotClause {
		query = query.Not(whereQuery, whereArgs)
	}

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	err := query.Find(&user).Error
	if err != nil {
		log.Errorf("[UserRepository][FindOne] failed to query First: %v", err)
		return nil, err
	}

	return user, err
}

func (u userRepository) FindOneWithAttribute(whereClause map[string]any, whereNotClause map[string]any, relations []string) (*models.User, error) {
	var user *models.User

	query := u.db

	for whereQuery, whereArgs := range whereClause {
		query = query.Where(whereQuery, whereArgs)
	}

	for whereQuery, whereArgs := range whereNotClause {
		query = query.Not(whereQuery, whereArgs)
	}

	for _, relation := range relations {
		query = query.Preload(relation)
	}

	err := query.First(&user).Error
	if err != nil {
		log.Errorf("[UserRepository][FindOne] failed to query First: %v", err)
		return nil, err
	}

	return user, err
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}
