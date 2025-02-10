package auth

import (
	"Webook/webook/internal/service/sms"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type AuthSMSService struct {
	svc sms.Service
	key []byte
}

type AuthSMSClaims struct {
	jwt.RegisteredClaims
	Tpl string
}

func (s *AuthSMSService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims AuthSMSClaims
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(t *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}

	return s.svc.Send(ctx, claims.Tpl, args, numbers...)
}
