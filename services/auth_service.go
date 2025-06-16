package services

import (
	"backend-attendance-deals/config"
	"backend-attendance-deals/database/models"
	"backend-attendance-deals/dto"
	"backend-attendance-deals/pkg/shared"
	"backend-attendance-deals/pkg/utils"
	"backend-attendance-deals/repositories"
	"context"
	"errors"
	"fmt"
	"strings"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, loginInput shared.LoginInput) (*shared.LoginResponse, error)
}

type authService struct {
	userRepository     repositories.UserRepositoryInterface
	logAuditRepository repositories.LogAuditRepositoryInterface
}

func (a authService) Login(ctx context.Context, loginInput shared.LoginInput) (*shared.LoginResponse, error) {
	ipAddress := utils.GetCurrentIPKey(&ctx)
	requestID := utils.GetCurrentRequestKey(&ctx)

	findUser, err := a.userRepository.FindOneWithAttribute(
		map[string]any{
			"LOWER(username) = ?": strings.ToLower(loginInput.Username),
		},
		nil,
		nil,
	)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	if findUser.Password != "" {
		err = findUser.ComparePasswords(loginInput.Password)

		if err != nil {
			return nil, errors.New("invalid username or password")
		}
	}

	accessTokenPayload := dto.AccessTokenPayload{
		UserId:   findUser.ID.String(),
		UserType: string(findUser.Role),
	}

	jwtAccessTokenExpiration := config.GetEnv("JWT_ACCESS_TOKEN_EXPIRATION", "24h")
	jwtSecret := config.GetEnv("JWT_SECRET", "")
	accessToken, err := utils.JwtGenerate(accessTokenPayload, jwtAccessTokenExpiration, jwtSecret)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to generate accessToken due to %s", err.Error()))
	}

	arrLogAudit := make([]*models.LogAudit, 0)
	arrLogAudit = append(arrLogAudit, &models.LogAudit{
		BaseModel: models.BaseModel{
			CreatedBy: &findUser.ID,
		},
		UserID:    findUser.ID,
		Action:    "login",
		Entity:    "User",
		EntityID:  findUser.ID,
		IpAddress: *ipAddress,
		RequestID: *requestID,
	})
	_, err = a.logAuditRepository.Create(ctx, arrLogAudit)

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return &shared.LoginResponse{
		AccessToken: *accessToken,
		User:        findUser,
	}, nil
}

func NewAuthService(
	userRepository repositories.UserRepositoryInterface,
	logAuditRepository repositories.LogAuditRepositoryInterface,
) AuthServiceInterface {
	return &authService{
		userRepository:     userRepository,
		logAuditRepository: logAuditRepository,
	}
}
