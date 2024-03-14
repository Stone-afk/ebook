package repository

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type userRepository struct {
	dao dao.UserDAO
}

func (repo *userRepository) InputUser(ctx context.Context, msg domain.User) error {
	return repo.dao.InputUser(ctx, dao.User{
		Id:       msg.Id,
		Email:    msg.Email,
		Nickname: msg.Nickname,
		Phone:    msg.Phone,
	})
}

func (repo *userRepository) SearchUser(ctx context.Context, keywords []string) ([]domain.User, error) {
	users, err := repo.dao.Search(ctx, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(users, func(idx int, src dao.User) domain.User {
		return domain.User{
			Id:       src.Id,
			Email:    src.Email,
			Nickname: src.Nickname,
			Phone:    src.Phone,
		}
	}), nil
}

func NewUserRepository(d dao.UserDAO) UserRepository {
	return &userRepository{
		dao: d,
	}
}
