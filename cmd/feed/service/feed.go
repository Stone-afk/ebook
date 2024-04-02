package service

import (
	"context"
	followv1 "ebook/cmd/api/proto/gen/followrelation/v1"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"golang.org/x/sync/errgroup"
	"sort"
	"sync"
)

type feedService struct {
	repo         repository.FeedEventRepo
	handlerMap   map[string]Handler
	followClient followv1.FollowServiceClient
}

func (svc *feedService) registerService(typ string, handler Handler) {
	svc.handlerMap[typ] = handler
}

func (svc *feedService) CreateFeedEvent(ctx context.Context, feed domain.FeedEvent) error {
	// 需要可以解决的handler
	handler, ok := svc.handlerMap[feed.Type]
	if !ok {
		// 这里你可以考虑引入一个兜底的处理机制。
		// 例如说在找不到的时候就默认丢过去 PushEvent 里面
		// 对于大部分业务来说，都是合适的
		return fmt.Errorf("未找到具体的业务处理逻辑 %s", feed.Type)
	}
	return handler.CreateFeedEvent(ctx, feed.Ext)
}

// GetFeedEventListV1 不依赖于 Handler 的直接查询
func (svc *feedService) GetFeedEventListV1(ctx context.Context, uid int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	var eg errgroup.Group
	var mu sync.RWMutex
	res := make([]domain.FeedEvent, 0, limit*2)
	eg.Go(func() error {
		resp, rerr := svc.followClient.GetFollowee(ctx, &followv1.GetFolloweeRequest{
			Follower: uid,
			Offset:   0,
			Limit:    200,
		})
		if rerr != nil {
			return rerr
		}
		followeeIds := slice.Map(resp.FollowRelations, func(idx int, src *followv1.FollowRelation) int64 {
			return src.Followee
		})
		events, err := svc.repo.FindPullEvents(ctx, followeeIds, timestamp, limit)
		if err != nil {
			return err
		}
		mu.Lock()
		res = append(res, events...)
		mu.Unlock()
		return nil
	})
	eg.Go(func() error {
		events, err := svc.repo.FindPushEvents(ctx, uid, timestamp, limit)
		if err != nil {
			return err
		}
		mu.Lock()
		res = append(res, events...)
		mu.Unlock()
		return nil
	})
	sort.Slice(res, func(i, j int) bool {
		return res[i].Ctime.Unix() > res[j].Ctime.Unix()
	})
	err := eg.Wait()
	return res[:slice.Min[int]([]int{int(limit), len(res)})], err
}

func (svc *feedService) GetFeedEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	var eg errgroup.Group
	var mu sync.RWMutex
	res := make([]domain.FeedEvent, 0, limit*int64(len(svc.handlerMap)))
	for _, handler := range svc.handlerMap {
		h := handler
		eg.Go(func() error {
			events, err := h.FindFeedEvents(ctx, uid, timestamp, limit)
			if err != nil {
				return err
			}
			mu.Lock()
			res = append(res, events...)
			mu.Unlock()
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	// 聚合排序，难免的
	sort.Slice(res, func(i, j int) bool {
		return res[i].Ctime.Unix() > res[j].Ctime.Unix()
	})
	return res[:slice.Min[int]([]int{int(limit), len(res)})], nil
}

func NewFeedService(repo repository.FeedEventRepo,
	//client followv1.FollowServiceClient,
	handlerMap map[string]Handler) FeedService {
	return &feedService{
		repo: repo,
		// 你可以注入那个 health client
		//followClient: client,
		handlerMap: handlerMap,
	}
}
