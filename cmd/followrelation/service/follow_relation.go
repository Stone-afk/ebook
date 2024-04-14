package service

import (
	"context"
	"ebook/cmd/followrelation/domain"
	"ebook/cmd/followrelation/events"
	"ebook/cmd/followrelation/repository"
	"ebook/cmd/pkg/logger"
	"strconv"
)

type followRelationService struct {
	feedEventProducer events.FeedEventProducer
	l                 logger.Logger
	repo              repository.FollowRepository
}

func (s *followRelationService) GetAllFollower(ctx context.Context, followee int64) ([]domain.FollowRelation, error) {
	offset := 0
	limit := 1000
	relations := make([]domain.FollowRelation, 0, limit)
	for {
		res, err := s.repo.GetFollower(ctx, followee, int64(offset), int64(limit))
		if err != nil {
			return nil, err
		}
		relations = append(relations, res...)
		if len(res) < limit {
			break
		}
		offset += limit
	}
	return relations, nil
}

func (s *followRelationService) sendFeedEvent(ctx context.Context, f domain.FollowRelation) {
	evt := events.FeedEvent{
		Type: "follower",
		Metadata: map[string]string{
			"followee": strconv.FormatInt(f.Followee, 10),
			"follower": strconv.FormatInt(f.Follower, 10),
			"biz_id":   strconv.FormatInt(f.Id, 10),
		},
	}
	er := s.feedEventProducer.ProduceFeedEvent(ctx, evt)
	if er != nil {
		s.l.Error("ProduceFeedEvent 发送feed流事件失败", logger.Error(er))
		er = s.feedEventProducer.ProduceStandardFeedEvent(ctx, evt)
		if er != nil {
			s.l.Error("ProduceStandardFeedEvent 发送feed流事件失败", logger.Error(er))
		}
	}
}

func (s *followRelationService) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	return s.repo.GetFollowStatics(ctx, uid)
}

func (s *followRelationService) GetFollower(ctx context.Context, followee, offset, limit int64) ([]domain.FollowRelation, error) {
	return s.repo.GetFollower(ctx, followee, offset, limit)
}

func (s *followRelationService) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	return s.repo.GetFollowee(ctx, follower, offset, limit)
}

func (s *followRelationService) FollowInfo(ctx context.Context, follower, followee int64) (domain.FollowRelation, error) {
	return s.repo.FollowInfo(ctx, follower, followee)
}

func (s *followRelationService) Follow(ctx context.Context, follower, followee int64) error {
	err := s.repo.AddFollowRelation(ctx, domain.FollowRelation{
		Followee: followee,
		Follower: follower,
	})
	if err != nil {
		return err
	}
	go func() {
		f, er := s.repo.FollowInfo(ctx, follower, followee)
		if er != nil {
			s.l.Error("获取关注详情信息失败", logger.Error(er))
			return
		}
		s.sendFeedEvent(ctx, f)
	}()
	return nil
}

func (s *followRelationService) CancelFollow(ctx context.Context, follower, followee int64) error {
	return s.repo.InactiveFollowRelation(ctx, follower, followee)
}

func NewFollowRelationService(feedEventProducer events.FeedEventProducer,
	l logger.Logger, repo repository.FollowRepository) FollowRelationService {
	return &followRelationService{
		feedEventProducer: feedEventProducer,
		repo:              repo,
		l:                 l,
	}
}
