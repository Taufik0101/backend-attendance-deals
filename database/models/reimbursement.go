package models

import (
	"github.com/google/uuid"
)

type Reimbursement struct {
	BaseModel
	Amount      float64           `json:"amount" gorm:"column:amount;"`
	Description string            `json:"description" gorm:"column:description;"`
	UserID      uuid.UUID         `json:"user_id" gorm:"column:user_id;"`
	User        *User             `json:"user,omitempty"`
	PeriodID    uuid.UUID         `json:"period_id" gorm:"column:period_id;"`
	Period      *AttendancePeriod `json:"period,omitempty"`
	IpAddress   string            `json:"ip_address" gorm:"ip_address"`
}

func (*Reimbursement) TableName() string {
	return "reimbursements"
}
