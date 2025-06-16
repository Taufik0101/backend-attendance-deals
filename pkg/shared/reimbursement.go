package shared

type (
	CreateReimbursementInput struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
		DayActivity string  `json:"day_activity"`
	}
)
