package utils

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtCustomClaim struct {
	Payload interface{} `json:"payload"`
	jwt.StandardClaims
}

func JwtGenerate(payload interface{}, duration string, secret string) (*string, error) {

	if secret == "" {
		return nil, fmt.Errorf("jwt secret cannot be empty")
	}

	durationParsed, err := time.ParseDuration(duration)
	if err != nil {
		log.Println("Duration error", err.Error())
		return nil, fmt.Errorf(fmt.Sprintf("failed to parse duration %s", err.Error()))
	}

	expiresAt := time.Now().Add(time.Duration(durationParsed)).Unix()
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  time.Now().Unix(),
		},
	})

	token, err := t.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func JwtValidate(token string, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return []byte(secret), nil
	})
}

func ExtractClaims(validatedToken jwt.Token) (jwt.MapClaims, error) {
	claims, ok := validatedToken.Claims.(jwt.MapClaims)
	if ok && validatedToken.Valid {
		return claims, nil
	}

	return nil, errors.New("failed to extract claims")
}
