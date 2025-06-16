package models

import (
	"github.com/google/uuid"
)

type Payroll struct {
	BaseModel
	PeriodID uuid.UUID         `json:"period_id" gorm:"column:period_id;"`
	Period   *AttendancePeriod `json:"period,omitempty"`
}

func (*Payroll) TableName() string {
	return "payrolls"
}
