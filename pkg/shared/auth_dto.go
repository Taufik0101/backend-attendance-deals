package shared

import "backend-attendance-deals/database/models"

type (
	LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		AccessToken string       `json:"accessToken"`
		User        *models.User `json:"user"`
	}
)
