package dao

import (
	"context"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"gorm.io/gorm"
)

const FollowRelationTableName = "follow_relations"

var (
	ErrFollowerNotFound = gorm.ErrRecordNotFound
)

type TableStoreFollowRelationDao struct {
	client *tablestore.TableStoreClient
}

func (dao *TableStoreFollowRelationDao) FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *TableStoreFollowRelationDao) FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *TableStoreFollowRelationDao) CreateFollowRelation(ctx context.Context, f FollowRelation) error {
	//TODO implement me
	panic("implement me")
}

func (dao *TableStoreFollowRelationDao) UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func (dao *TableStoreFollowRelationDao) CntFollower(ctx context.Context, uid int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (dao *TableStoreFollowRelationDao) CntFollowee(ctx context.Context, uid int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func NewTableStoreDao(client *tablestore.TableStoreClient) FollowRelationDao {
	return &TableStoreFollowRelationDao{
		client: client,
	}
}
