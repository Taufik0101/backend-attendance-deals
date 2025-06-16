package dto

import "github.com/google/uuid"

type AuthToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenPayload struct {
	UserId   string `json:"user_id" mapstructure:"user_id"`
	UserType string `json:"user_type" mapstructure:"user_type"`
}

type RefreshTokenPayload struct {
	UserId string `json:"user_id"`
}

func (d AccessTokenPayload) UserIDStrToUUID() uuid.UUID {
	id, _ := uuid.Parse(d.UserId)
	return id
}
