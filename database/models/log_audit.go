package models

import (
	"github.com/google/uuid"
)

type LogAudit struct {
	BaseModel
	UserID    uuid.UUID         `json:"user_id" gorm:"column:user_id;"`
	User      *User             `json:"user,omitempty"`
	PeriodID  *uuid.UUID        `json:"period_id" gorm:"column:period_id;"`
	Period    *AttendancePeriod `json:"period,omitempty"`
	Action    string            `json:"action" gorm:"column:action"`
	Entity    string            `json:"entity" gorm:"column:entity"`
	EntityID  uuid.UUID         `json:"entity_id" gorm:"column:entity_id"`
	IpAddress string            `json:"ip_address" gorm:"ip_address"`
	RequestID string            `json:"request_id" gorm:"request_id"`
}

func (*LogAudit) TableName() string {
	return "audit_logs"
}
