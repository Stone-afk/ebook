package service

import (
	"context"
	"ebook/cmd/internal/repository"
	"ebook/cmd/pkg/logger"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/service/interactive.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/service/mocks/interactive.mock.go
type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
}

type interactiveService struct {
	repo repository.InteractiveRepository
	l    logger.Logger
}

func (svc *interactiveService) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	// 点赞
	return svc.repo.IncrLike(ctx, biz, bizId, uid)
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	return svc.repo.DecrLike(ctx, biz, bizId, uid)
}