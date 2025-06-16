package utils

import (
	"backend-attendance-deals/dto"
	"context"
	"github.com/gofiber/fiber/v2"
)

type KeyType string

const (
	AuthSuccessKey      KeyType = "authSuccess"
	CurrentTimezoneKey  KeyType = "currentTimezone"
	CurrentUserKey      KeyType = "currentUser"
	CurrentRefererKey   KeyType = "currentReferer"
	CurrentUserRolesKey KeyType = "currentUserRolesKey"
	ErrorKey            KeyType = "error"
	FiberContextKey     KeyType = "fiberContextKey"
	CurrentRequestKey   KeyType = "currentRequestKey"
	CurrentIPKey        KeyType = "currentIPKey"
)

func GetCurrentTimezone(ctx *context.Context) *string {
	timezone := (*ctx).Value(CurrentTimezoneKey)
	if timezone != nil {
		tz := timezone.(string)
		return &tz
	}

	return nil
}

func GetCurrentUser(ctx *context.Context) *dto.AccessTokenPayload {
	fiberCtx := (*ctx).Value("fiberContext").(*fiber.Ctx)
	user := fiberCtx.Locals(CurrentUserKey)
	if user != nil {
		return user.(*dto.AccessTokenPayload)
	}

	return nil
}

func GetCurrentRefererKey(ctx *context.Context) *string {
	fiberCtx := (*ctx).Value("fiberContext").(*fiber.Ctx)
	referer := fiberCtx.Locals(CurrentRefererKey)
	if referer != nil {
		return referer.(*string)
	}

	return nil
}

func GetCurrentRequestKey(ctx *context.Context) *string {
	fiberCtx := (*ctx).Value("fiberContext").(*fiber.Ctx)
	referer := fiberCtx.Locals(CurrentRequestKey)
	if referer != nil {
		ref := referer.(string)
		return &ref
	}

	return nil
}

func GetCurrentIPKey(ctx *context.Context) *string {
	fiberCtx := (*ctx).Value("fiberContext").(*fiber.Ctx)
	referer := fiberCtx.Locals(CurrentIPKey)
	if referer != nil {
		ref := referer.(string)
		return &ref
	}

	return nil
}
