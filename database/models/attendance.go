package models

import (
	"github.com/google/uuid"
	"time"
)

type Attendance struct {
	BaseModel
	DateIn    time.Time         `json:"date_in" gorm:"column:date_in;"`
	DateOut   *time.Time        `json:"date_out" gorm:"column:date_out;"`
	UserID    uuid.UUID         `json:"user_id" gorm:"column:user_id;"`
	User      *User             `json:"user,omitempty"`
	PeriodID  uuid.UUID         `json:"period_id" gorm:"column:period_id;"`
	Period    *AttendancePeriod `json:"period,omitempty"`
	IpAddress string            `json:"ip_address" gorm:"ip_address"`
}

func (*Attendance) TableName() string {
	return "attendances"
}
