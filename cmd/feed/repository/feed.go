package repository

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository/cache"
	"ebook/cmd/feed/repository/dao"
)

type feedEventRepo struct {
	pullDao   dao.FeedPullEventDAO
	pushDao   dao.FeedPushEventDAO
	feedCache cache.FeedEventCache
}

func (repo *feedEventRepo) CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error {
	//TODO implement me
	panic("implement me")
}

func (repo *feedEventRepo) CreatePullEvent(ctx context.Context, event domain.FeedEvent) error {
	//TODO implement me
	panic("implement me")
}

func (repo *feedEventRepo) FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *feedEventRepo) FindPushEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *feedEventRepo) FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *feedEventRepo) FindPushEventsWithTyp(ctx context.Context, typ string, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	//TODO implement me
	panic("implement me")
}

func NewFeedEventRepo(pullDao dao.FeedPullEventDAO, pushDao dao.FeedPushEventDAO, feedCache cache.FeedEventCache) FeedEventRepo {
	return &feedEventRepo{
		pullDao:   pullDao,
		pushDao:   pushDao,
		feedCache: feedCache,
	}
}
