package models

import (
	"time"
)

type AttendancePeriod struct {
	BaseModel
	StartDate time.Time `json:"start_date" gorm:"column:start_date;"`
	EndDate   time.Time `json:"end_date" gorm:"column:end_date;"`
}

func (*AttendancePeriod) TableName() string {
	return "attendance_periods"
}
