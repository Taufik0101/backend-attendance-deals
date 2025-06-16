package models

import (
	"github.com/google/uuid"
	"time"
)

type Overtime struct {
	BaseModel
	Date      time.Time         `json:"date" gorm:"column:date;"`
	Hours     int               `json:"hours" gorm:"column:hours;"`
	UserID    uuid.UUID         `json:"user_id" gorm:"column:user_id;"`
	User      *User             `json:"user,omitempty"`
	PeriodID  uuid.UUID         `json:"period_id" gorm:"column:period_id;"`
	Period    *AttendancePeriod `json:"period,omitempty"`
	IpAddress string            `json:"ip_address" gorm:"ip_address"`
}

func (*Overtime) TableName() string {
	return "overtimes"
}
