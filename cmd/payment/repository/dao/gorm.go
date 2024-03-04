package dao

import (
	"context"
	"ebook/cmd/payment/domain"
	"gorm.io/gorm"
	"time"
)

type PaymentGORMDAO struct {
	db *gorm.DB
}

func (d PaymentGORMDAO) Insert(ctx context.Context, pmt Payment) error {
	now := time.Now().UnixMilli()
	pmt.Utime = now
	pmt.Ctime = now
	return d.db.WithContext(ctx).Create(&pmt).Error
}

func (d PaymentGORMDAO) UpdateTxnIDAndStatus(ctx context.Context, bizTradeNo string, txnID string, status domain.PaymentStatus) error {
	//TODO implement me
	panic("implement me")
}

func (d PaymentGORMDAO) FindExpiredPayment(ctx context.Context, offset int, limit int, t time.Time) ([]Payment, error) {
	//TODO implement me
	panic("implement me")
}

func (d PaymentGORMDAO) GetPayment(ctx context.Context, bizTradeNO string) (Payment, error) {
	var res Payment
	err := d.db.WithContext(ctx).Where("biz_trade_no = ?", bizTradeNO).First(&res).Error
	return res, err
}

func NewPaymentGORMDAO(db *gorm.DB) PaymentDAO {
	return &PaymentGORMDAO{db: db}
}
