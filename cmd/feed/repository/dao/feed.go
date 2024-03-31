package dao

import (
	"context"
	"gorm.io/gorm"
)

type feedPullEventDAO struct {
	db *gorm.DB
}

func (dao *feedPullEventDAO) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	//TODO implement me
	panic("implement me")
}

func (dao *feedPullEventDAO) FindPullEventList(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *feedPullEventDAO) FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	//TODO implement me
	panic("implement me")
}

func NewFeedPullEventDAO(db *gorm.DB) FeedPullEventDAO {
	return &feedPullEventDAO{
		db: db,
	}
}

type feedPushEventDAO struct {
	db *gorm.DB
}

func (dao *feedPushEventDAO) CreatePushEvents(ctx context.Context, events []FeedPushEvent) error {
	//TODO implement me
	panic("implement me")
}

func (dao *feedPushEventDAO) GetPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *feedPushEventDAO) GetPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	//TODO implement me
	panic("implement me")
}

func NewFeedPushEventDAO(db *gorm.DB) FeedPushEventDAO {
	return &feedPushEventDAO{
		db: db,
	}
}
