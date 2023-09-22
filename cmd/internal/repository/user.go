package repository

import (
	"context"
	"database/sql"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao"
	"time"
)

var ErrUserDuplicate = dao.ErrUserDuplicate
var ErrUserNotFound = dao.ErrDataNotFound

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/user.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/user.mock.go

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByWechat(ctx context.Context, openID string) (domain.User, error)

	// Update 更新数据，只有非 0 值才会更新
	Update(ctx context.Context, u domain.User) error
}

// UserRepository 使用了缓存的 repository 实现
type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

// NewUserRepository 也说明了 CachedUserRepository 的特性
// 会从缓存和数据库中去尝试获得
func NewUserRepository(d dao.UserDAO,
	c cache.UserCache) UserRepository {
	return &userRepository{
		dao:   d,
		cache: c,
	}
}

func (ur *userRepository) Update(ctx context.Context, u domain.User) error {
	err := ur.dao.UpdateNonZeroFields(ctx, ur.domainToEntity(u))
	if err != nil {
		return err
	}
	return ur.cache.Delete(ctx, u.Id)
}

func (ur *userRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domainToEntity(u))
}

func (ur *userRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := ur.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(u), err
}

func (ur *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	return ur.entityToDomain(u), err
}

func (ur *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	return ur.entityToDomain(u), err
}

func (ur *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := ur.cache.Get(ctx, id)
	if err == nil {
		return u, err
	}
	ue, err := ur.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = ur.entityToDomain(ue)
	_ = ur.cache.Set(ctx, u)
	return u, nil
}

func (ur *userRepository) entityToDomain(ue dao.User) domain.User {
	var birthday time.Time
	if ue.Birthday.Valid {
		birthday = time.UnixMilli(ue.Birthday.Int64)
	}
	return domain.User{
		Id:       ue.Id,
		Email:    ue.Email.String,
		Password: ue.Password,
		Phone:    ue.Phone.String,
		Nickname: ue.Nickname.String,
		AboutMe:  ue.AboutMe.String,
		Birthday: birthday,
		UnionID:  ue.WechatUnionID.String,
		OpenID:   ue.WechatOpenID.String,
		Ctime:    time.UnixMilli(ue.Ctime),
	}
}

func (ur *userRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Birthday: sql.NullInt64{
			Int64: u.Birthday.UnixMilli(),
			Valid: u.Birthday.IsZero(),
		},
		Nickname: sql.NullString{
			String: u.Nickname,
			Valid:  u.Nickname != "",
		},
		AboutMe: sql.NullString{
			String: u.AboutMe,
			Valid:  u.AboutMe != "",
		},
		Password: u.Password,
		WechatOpenID: sql.NullString{
			String: u.OpenID,
			Valid:  u.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.UnionID,
			Valid:  u.UnionID != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}
