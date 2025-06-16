package shared

import "time"

type (
	CreatePayslipInput struct {
		PeriodID string `json:"period_id"`
	}

	ListSummaryPaySlip struct {
		PeriodID string `json:"period_id"`
	}

	DetailPaySlipResponse struct {
		User        string  `json:"user"`
		TakeHomePay float64 `json:"take_home_pay"`
	}

	PaySlipResponse struct {
		TotalTakeHomePay  float64                 `json:"total_take_home_pay"`
		DetailTakeHomePay []DetailPaySlipResponse `json:"detail_take_home_pay"`
	}

	DetailOvertimes struct {
		Day   time.Time `json:"day"`
		Hours int       `json:"hours"`
	}

	Reimbursement struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	CreatePayslipResponse struct {
		Period struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		} `json:"period"`
		Attendances struct {
			TotalShouldAttend int         `json:"total_should_attend"`
			SalaryDay         float64     `json:"salary_day"`
			TotalAttend       int         `json:"total_attend"`
			Formula           string      `json:"formula"`
			DetailAttend      []time.Time `json:"detail_attend"`
		} `json:"attendances"`
		Overtimes struct {
			SalaryHours     float64           `json:"salary_hours"`
			Formula         string            `json:"formula"`
			OvertimeSalary  float64           `json:"overtime_salary"`
			TotalOvertime   int               `json:"total_overtime"`
			DetailOvertimes []DetailOvertimes `json:"detail_overtimes"`
		} `json:"overtimes"`
		Reimbursements []Reimbursement `json:"reimbursements"`
		TakeHomePay    float64         `json:"take_home_pay"`
	}
)
