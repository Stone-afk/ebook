package repository

import (
	"context"
	"ebook/cmd/payment/domain"
	"ebook/cmd/payment/repository/dao"
	"time"
)

type paymentRepository struct {
	dao dao.PaymentDAO
}

func (p paymentRepository) AddPayment(ctx context.Context, pmt domain.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (p paymentRepository) UpdatePayment(ctx context.Context, pmt domain.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (p paymentRepository) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (p paymentRepository) GetPayment(ctx context.Context, bizTradeNO string) (domain.Payment, error) {
	//TODO implement me
	panic("implement me")
}

func NewPaymentRepository(d dao.PaymentDAO) PaymentRepository {
	return &paymentRepository{
		dao: d,
	}
}
