package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type TagElasticDAO struct {
	client *elastic.Client
}

func (dao *TagElasticDAO) InputTag(ctx context.Context, tag Tag) error {
	//TODO implement me
	panic("implement me")
}

func (dao *TagElasticDAO) Search(ctx context.Context, uid int64, biz string, keywords []string) ([]Tag, error) {
	//TODO implement me
	panic("implement me")
}

func NewTagElasticDAO(client *elastic.Client) TagDAO {
	return &TagElasticDAO{client: client}
}
