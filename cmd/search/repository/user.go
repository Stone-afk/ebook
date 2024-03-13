package repository

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository/dao"
)

type userRepository struct {
	dao dao.UserDAO
}

func (dao *userRepository) InputUser(ctx context.Context, msg domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (dao *userRepository) SearchUser(ctx context.Context, keywords []string) ([]domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(d dao.UserDAO) UserRepository {
	return &userRepository{
		dao: d,
	}
}
