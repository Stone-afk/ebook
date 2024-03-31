package repository

import (
	"context"
	"ebook/cmd/feed/domain"
	"ebook/cmd/feed/repository/cache"
	"ebook/cmd/feed/repository/dao"
	"encoding/json"
	"time"
)

type feedEventRepo struct {
	pullDao   dao.FeedPullEventDAO
	pushDao   dao.FeedPushEventDAO
	feedCache cache.FeedEventCache
}

func (repo *feedEventRepo) CreatePushEvents(ctx context.Context, events []domain.FeedEvent) error {
	pushEvents := make([]dao.FeedPushEvent, 0, len(events))
	for _, e := range events {
		pushEvents = append(pushEvents, convertToPushEventDao(e))
	}
	return repo.pushDao.CreatePushEvents(ctx, pushEvents)
}

func (repo *feedEventRepo) CreatePullEvent(ctx context.Context, event domain.FeedEvent) error {
	return repo.pullDao.CreatePullEvent(ctx, convertToPullEventDao(event))
}

func (repo *feedEventRepo) FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := repo.pullDao.FindPullEvents(ctx, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPullEventDomain(e))
	}
	return ans, nil
}

func (repo *feedEventRepo) FindPushEvents(ctx context.Context, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := repo.pushDao.FindPushEvents(ctx, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPushEventDomain(e))
	}
	return ans, nil
}

func (repo *feedEventRepo) FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := repo.pullDao.FindPullEventsWithTyp(ctx, typ, uids, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPullEventDomain(e))
	}
	return ans, nil
}

func (repo *feedEventRepo) FindPushEventsWithTyp(ctx context.Context, typ string, uid, timestamp, limit int64) ([]domain.FeedEvent, error) {
	events, err := repo.pushDao.FindPushEventsWithTyp(ctx, typ, uid, timestamp, limit)
	if err != nil {
		return nil, err
	}
	ans := make([]domain.FeedEvent, 0, len(events))
	for _, e := range events {
		ans = append(ans, convertToPushEventDomain(e))
	}
	return ans, nil
}

func NewFeedEventRepo(pullDao dao.FeedPullEventDAO, pushDao dao.FeedPushEventDAO, feedCache cache.FeedEventCache) FeedEventRepo {
	return &feedEventRepo{
		pullDao:   pullDao,
		pushDao:   pushDao,
		feedCache: feedCache,
	}
}

func convertToPushEventDomain(event dao.FeedPushEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		ID:    event.Id,
		Uid:   event.UID,
		Type:  event.Type,
		Ctime: time.Unix(event.Ctime, 0),
		Ext:   ext,
	}
}

func convertToPullEventDomain(event dao.FeedPullEvent) domain.FeedEvent {
	var ext map[string]string
	_ = json.Unmarshal([]byte(event.Content), &ext)
	return domain.FeedEvent{
		ID:    event.Id,
		Uid:   event.UID,
		Type:  event.Type,
		Ctime: time.Unix(event.Ctime, 0),
		Ext:   ext,
	}
}

func convertToPushEventDao(event domain.FeedEvent) dao.FeedPushEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPushEvent{
		Id:      event.ID,
		UID:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.Unix(),
	}
}

func convertToPullEventDao(event domain.FeedEvent) dao.FeedPullEvent {
	val, _ := json.Marshal(event.Ext)
	return dao.FeedPullEvent{
		Id:      event.ID,
		UID:     event.Uid,
		Type:    event.Type,
		Content: string(val),
		Ctime:   event.Ctime.Unix(),
	}
}
