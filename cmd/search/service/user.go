package service

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository"
	"strings"
)

type userSearchService struct {
	userRepo repository.UserRepository
}

func (s *searchService) SearchUser(ctx context.Context, expression string) (domain.SearchUserResult, error) {
	keywords := strings.Split(expression, " ")
	var res domain.SearchUserResult
	users, err := s.userRepo.SearchUser(ctx, keywords)
	res.Users = users
	return res, err
}

func NewUserSearchService(userRepo repository.UserRepository) UserSearchService {
	return &searchService{userRepo: userRepo}
}
