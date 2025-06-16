package models

import (
	"github.com/google/uuid"
)

type Payslip struct {
	BaseModel
	AttendanceDays int       `json:"attendance_days" gorm:"column:attendance_days;"`
	OvertimeHours  int       `json:"overtime_hours" gorm:"column:overtime_hours;"`
	Reimbursement  float64   `json:"reimbursement" gorm:"column:reimbursement;"`
	BaseSalary     float64   `json:"base_salary" gorm:"column:base_salary;"`
	TotalPay       float64   `json:"total_pay" gorm:"column:total_pay;"`
	UserID         uuid.UUID `json:"user_id" gorm:"column:user_id;"`
	User           *User     `json:"user,omitempty"`
	PayrollID      uuid.UUID `json:"payroll_id" gorm:"column:payroll_id;"`
	Payroll        *Payroll  `json:"payroll,omitempty"`
}

func (*Payslip) TableName() string {
	return "payslips"
}
