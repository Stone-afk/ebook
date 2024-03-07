package repository

import (
	"context"
	"ebook/cmd/account/domain"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
