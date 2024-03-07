package service

import (
	"context"
	"ebook/cmd/account/domain"
	"ebook/cmd/account/repository"
)

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (svc *accountService) Credit(ctx context.Context, cr domain.Credit) error {
	return svc.repo.AddCredit(ctx, cr)
}
