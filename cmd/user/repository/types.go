package repository

import (
	"context"
	"ebook/cmd/user/domain"
	"ebook/cmd/user/repository/dao"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrDataNotFound

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/user/repository/types.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/user/repository/mocks/user.mock.go
type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
	// Update 更新数据，只有非 0 值才会更新
	Update(ctx context.Context, u domain.User) error
}

