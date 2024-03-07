package repository

import (
	"context"
	"ebook/cmd/account/domain"
	"ebook/cmd/account/repository/dao"
	"time"
)

type accountRepository struct {
	dao dao.AccountDAO
}

func (repo *accountRepository) AddCredit(ctx context.Context, c domain.Credit) error {
	activities := make([]dao.AccountActivity, 0, len(c.Items))
	now := time.Now().UnixMilli()
	for _, itm := range c.Items {
		activities = append(activities, dao.AccountActivity{
			Uid:         itm.Uid,
			Biz:         c.Biz,
			BizId:       c.BizId,
			Account:     itm.Account,
			AccountType: itm.AccountType.AsUint8(),
			Amount:      itm.Amt,
			Currency:    itm.Currency,
			Ctime:       now,
			Utime:       now,
		})
	}
	return repo.dao.AddActivities(ctx, activities...)
}

func NewAccountRepository(dao dao.AccountDAO) AccountRepository {
	return &accountRepository{dao: dao}
}
