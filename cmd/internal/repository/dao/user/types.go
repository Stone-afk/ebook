package user

import (
	"context"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/user/user.go -package=daomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/dao/mocks/user.mock.go
type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByWechat(ctx context.Context, openID string) (User, error)
	UpdateNonZeroFields(ctx context.Context, u User) error
}
