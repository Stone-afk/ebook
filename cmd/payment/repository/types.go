package repository

import (
	"context"
	"ebook/cmd/payment/domain"
	"time"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/payment/repository/payment.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/payment/repository/mocks/payment.mock.go
type PaymentRepository interface {
	AddPayment(ctx context.Context, pmt domain.Payment) error
	// UpdatePayment 这个设计有点差，因为
	UpdatePayment(ctx context.Context, pmt domain.Payment) error
	FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error)
	GetPayment(ctx context.Context, bizTradeNO string) (domain.Payment, error)
}
