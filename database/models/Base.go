package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	/**
	notes: if migrations fail because uuid_generate_v4()
	is not exists please run code below
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; on postgres console
	*/
	ID         uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreatedBy  *uuid.UUID     `json:"created_by" gorm:"column:created_by;type:uuid;default:null"`
	UpdatedBy  *uuid.UUID     `json:"updated_by" gorm:"column:updated_by;type:uuid;default:null"`
	CreatedBys *User          `json:"created_bys,omitempty" gorm:"foreignKey:ID;references:CreatedBy"`
	UpdatedBys *User          `json:"updated_bys,omitempty" gorm:"foreignKey:ID;references:UpdatedBy"`
	CreatedAt  time.Time      `json:"created_at" gorm:"not null;type:timestamp;autoCreateTime;column:created_at"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"not null;type:timestamp;autoUpdateTime;column:updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (b *BaseModel) Get() (interface{}, error) {
	return nil, nil
}
