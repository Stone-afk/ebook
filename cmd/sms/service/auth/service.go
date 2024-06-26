package auth

import (
	"context"
	"ebook/cmd/sms/service"
	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	svc service.Service
	key []byte
}

func NewService(svc service.Service, key []byte) *Service {
	return &Service{
		svc: svc,
		key: key,
	}
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}

func (s *Service) Send(ctx context.Context,
	tplToken string, args []string, numbers ...string) error {
	var c Claims
	_, err := jwt.ParseWithClaims(tplToken, &c, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, c.Tpl, args, numbers...)
}
