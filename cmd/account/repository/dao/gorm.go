package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type AccountGORMDAO struct {
	db *gorm.DB
}

func (dao *AccountGORMDAO) AddActivities(ctx context.Context, activities ...AccountActivity) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		for _, act := range activities {
			err := tx.Clauses(clause.OnConflict{
				DoUpdates: clause.Assignments(map[string]interface{}{
					"balance": gorm.Expr("balance + ?", act.Amount),
					"utime":   now,
				}),
			}).Create(&Account{
				Uid:      act.Uid,
				Account:  act.Account,
				Type:     act.AccountType,
				Balance:  act.Amount,
				Currency: act.Currency,
				Ctime:    now,
				Utime:    now,
			}).Error
			if err != nil {
				return err
			}
		}
		return tx.Create(&activities).Error
	})
}

func NewCreditGORMDAO(db *gorm.DB) AccountDAO {
	return &AccountGORMDAO{db: db}
}
