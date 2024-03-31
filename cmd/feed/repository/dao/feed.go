package dao

import (
	"context"
	"gorm.io/gorm"
)

type feedPullEventDAO struct {
	db *gorm.DB
}

func (dao *feedPullEventDAO) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return dao.db.WithContext(ctx).Create(&event).Error
}

func (dao *feedPullEventDAO) FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := dao.db.WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (dao *feedPullEventDAO) FindPullEventsWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := dao.db.WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
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
	return dao.db.WithContext(ctx).Create(events).Error
}

func (dao *feedPushEventDAO) FindPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := dao.db.WithContext(ctx).
		Where("uid = ?", uid).
		Where("ctime < ?", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (dao *feedPushEventDAO) FindPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := dao.db.WithContext(ctx).
		Where("uid = ?", uid).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func NewFeedPushEventDAO(db *gorm.DB) FeedPushEventDAO {
	return &feedPushEventDAO{
		db: db,
	}
}
