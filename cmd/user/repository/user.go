package repository

import (
	"context"
	"database/sql"
	"ebook/cmd/user/domain"
	"ebook/cmd/user/repository/cache"
	"ebook/cmd/user/repository/dao"
	"time"
)

// CachedUserRepository 使用了缓存的 repository 实现
type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

// NewUserRepository 也说明了 CachedUserRepository 的特性
// 会从缓存和数据库中去尝试获得
func NewUserRepository(d dao.UserDAO,
	c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   d,
		cache: c,
	}
}

func (ur *CachedUserRepository) Update(ctx context.Context, u domain.User) error {
	err := ur.dao.UpdateNonZeroFields(ctx, ur.domainToEntity(u))
	if err != nil {
		return err
	}
	return ur.cache.Delete(ctx, u.Id)
}

func (ur *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.domainToEntity(u))
}

func (ur *CachedUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := ur.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return ur.entityToDomain(u), err
}

func (ur *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := ur.dao.FindByPhone(ctx, phone)
	return ur.entityToDomain(u), err
}

func (ur *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	return ur.entityToDomain(u), err
}

func (ur *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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

func (ur *CachedUserRepository) entityToDomain(ue dao.User) domain.User {
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
		WechatInfo: domain.WechatInfo{
			// OpenId 是应用内唯一
			OpenId: ue.WechatOpenID.String,
			// UnionId 是整个公司账号内唯一
			UnionId: ue.WechatUnionID.String,
		},
		Ctime: time.UnixMilli(ue.Ctime),
	}
}

func (ur *CachedUserRepository) domainToEntity(u domain.User) dao.User {
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
			String: u.WechatInfo.OpenId,
			Valid:  u.WechatInfo.OpenId != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionId,
			Valid:  u.WechatInfo.UnionId != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}
