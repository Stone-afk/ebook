package service

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository"
)

type UserService interface {
	Signup(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx context.Context, u domain.User) error {
	panic("")
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	panic("")
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	panic("")
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	panic("")
}
