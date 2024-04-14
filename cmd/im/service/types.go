package service

import (
	"context"
	"ebook/cmd/im/domain"
)

type Secret string

type BaseHost string

type UserService interface {
	Sync(ctx context.Context, user domain.User) error
}
