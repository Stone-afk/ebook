package service

import (
	"context"
	"ebook/cmd/user/domain"
	"ebook/cmd/user/repository"
	"errors"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("邮箱或者密码不正确")

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/user/service/types.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/user/service/mocks/user.mock.go
type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	// UpdateNonSensitiveInfo 更新非敏感数据
	// 你可以在这里进一步补充究竟哪些数据会被更新
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
}
