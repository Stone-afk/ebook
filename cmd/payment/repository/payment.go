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

func (repo *paymentRepository) AddPayment(ctx context.Context, pmt domain.Payment) error {
	return repo.dao.Insert(ctx, repo.toEntity(pmt))
}

func (repo *paymentRepository) UpdatePayment(ctx context.Context, pmt domain.Payment) error {
	return repo.dao.UpdateTxnIDAndStatus(ctx, pmt.BizTradeNO, pmt.TxnID, pmt.Status)
}

func (repo *paymentRepository) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]domain.Payment, error) {
	pmts, err := repo.dao.FindExpiredPayment(ctx, offset, limit, t)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Payment, 0, len(pmts))
	for _, pmt := range pmts {
		res = append(res, repo.toDomain(pmt))
	}
	return res, nil
}

func (repo *paymentRepository) GetPayment(ctx context.Context, bizTradeNO string) (domain.Payment, error) {
	r, err := repo.dao.GetPayment(ctx, bizTradeNO)
	return repo.toDomain(r), err
}

func NewPaymentRepository(d dao.PaymentDAO) PaymentRepository {
	return &paymentRepository{
		dao: d,
	}
}

func (repo *paymentRepository) toDomain(pmt dao.Payment) domain.Payment {
	return domain.Payment{
		Amt: domain.Amount{
			Currency: pmt.Currency,
			Total:    pmt.Amt,
		},
		BizTradeNO:  pmt.BizTradeNO,
		Description: pmt.Description,
		Status:      domain.PaymentStatus(pmt.Status),
		TxnID:       pmt.TxnID.String,
	}
}

func (repo *paymentRepository) toEntity(pmt domain.Payment) dao.Payment {
	return dao.Payment{
		Amt:         pmt.Amt.Total,
		Currency:    pmt.Amt.Currency,
		BizTradeNO:  pmt.BizTradeNO,
		Description: pmt.Description,
		Status:      domain.PaymentStatusInit,
	}
}
