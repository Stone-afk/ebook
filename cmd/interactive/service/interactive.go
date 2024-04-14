package service

import (
	"context"
	"ebook/cmd/interactive/domain"
	"ebook/cmd/interactive/events/like"
	"ebook/cmd/interactive/repository"
	"ebook/cmd/pkg/logger"
	"golang.org/x/sync/errgroup"
	"strconv"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/interactive/service/interactive.go -package=svcmocks -destination=/Users/stone/go_project/ebook/ebook/cmd/interactive/service/mocks/interactive.mock.go
type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	// Like 点赞
	Like(ctx context.Context, biz string, bizId int64, userId int64) error
	// CancelLike 取消点赞
	CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error
	Get(ctx context.Context, biz string, bizId, userId int64) (domain.Interactive, error)
	// Collect 收藏
	Collect(ctx context.Context, biz string, bizId, cid, uid int64) error
	GetByIds(ctx context.Context, biz string, bizIds []int64) (map[int64]domain.Interactive, error)
}

type interactiveService struct {
	feedEventProducer like.FeedEventProducer
	repo              repository.InteractiveRepository
	l                 logger.Logger
}

func NewInteractiveService(feedEventProducer like.FeedEventProducer,
	repo repository.InteractiveRepository,
	l logger.Logger) InteractiveService {
	return &interactiveService{
		feedEventProducer: feedEventProducer,
		repo:              repo,
		l:                 l,
	}
}

func (svc *interactiveService) GetByIds(ctx context.Context, biz string,
	bizIds []int64) (map[int64]domain.Interactive, error) {
	intrs, err := svc.repo.GetByIds(ctx, biz, bizIds)
	if err != nil {
		return nil, err
	}
	res := make(map[int64]domain.Interactive, len(intrs))
	for _, intr := range intrs {
		res[intr.BizId] = intr
	}
	return res, nil
}

func (svc *interactiveService) Get(ctx context.Context, biz string, bizId, userId int64) (domain.Interactive, error) {
	// 按照 repository 的语义(完成 domain.Interactive 的完整构造)，你这里拿到的就应该是包含全部字段的
	var (
		eg        errgroup.Group
		intr      domain.Interactive
		liked     bool
		collected bool
	)
	eg.Go(func() error {
		var err error
		intr, err = svc.repo.Get(ctx, biz, bizId)
		return err
	})
	eg.Go(func() error {
		var err error
		liked, err = svc.repo.Liked(ctx, biz, bizId, userId)
		return err
	})
	eg.Go(func() error {
		var err error
		collected, err = svc.repo.Collected(ctx, biz, bizId, userId)
		return err
	})
	err := eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	intr.Liked = liked
	intr.Collected = collected
	return intr, err
}

func (svc *interactiveService) sendFeedEvent(ctx context.Context, userId int64, inter domain.Interactive) {
	var liked int
	if inter.Liked {
		liked = 1
	}
	evt := like.FeedEvent{
		Type: "follower",
		Metadata: map[string]string{
			"liked":  strconv.Itoa(liked),
			"biz":    inter.Biz,
			"biz_id": strconv.FormatInt(inter.BizId, 10),
			"uid":    strconv.FormatInt(userId, 10),
		},
	}
	er := svc.feedEventProducer.ProduceFeedEvent(ctx, evt)
	if er != nil {
		svc.l.Error("ProduceFeedEvent 发送feed流事件失败", logger.Error(er))
		er = svc.feedEventProducer.ProduceStandardFeedEvent(ctx, evt)
		if er != nil {
			svc.l.Error("ProduceStandardFeedEvent 发送feed流事件失败", logger.Error(er))
		}
	}
}

func (svc *interactiveService) Like(ctx context.Context, biz string, bizId int64, userId int64) error {
	// 点赞
	err := svc.repo.IncrLike(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	go func() {
		inter, er := svc.Get(ctx, biz, bizId, userId)
		if er != nil {
			svc.l.Error("获取喜欢详情信息失败", logger.Error(er))
			return
		}
		svc.sendFeedEvent(ctx, userId, inter)
	}()
	return nil
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}

func (svc *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error {
	return svc.repo.DecrLike(ctx, biz, bizId, userId)
}

// Collect 收藏
func (svc *interactiveService) Collect(ctx context.Context,
	biz string, bizId, cid, uid int64) error {
	return svc.repo.AddCollectionItem(ctx, biz, bizId, cid, uid)
}
