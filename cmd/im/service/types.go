package service

import (
	"context"
	"ebook/cmd/im/domain"
)

type UserService interface {
	Sync(ctx context.Context, user domain.User) error
}
