package dao

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
)

// ErrUserDuplicate 这个算是 user 专属的
var ErrUserDuplicate = errors.New("用户邮箱或者手机号冲突")

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewGORMUserDAO(db *gorm.DB) *GORMUserDAO {
	return &GORMUserDAO{
		db: db,
	}
}

func (ud *GORMUserDAO) Insert(ctx context.Context, u User) error {
	panic("")
}

func (ud *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	panic("")
}

func (ud *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	panic("")
}

func (ud *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	panic("")
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置为唯一索引
	Email    sql.NullString `gorm:"unique"`
	Password string

	//Phone *string
	Phone sql.NullString `gorm:"unique"`

	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
