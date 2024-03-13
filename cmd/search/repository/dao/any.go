package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type AnyElasticDAO struct {
	client *elastic.Client
}

func NewAnyElasticDAO(client *elastic.Client) AnyDAO {
	return &AnyElasticDAO{client: client}
}

func (dao *AnyElasticDAO) Input(ctx context.Context, index, docId, data string) error {
	_, err := dao.client.Index().
		Index(index).Id(docId).BodyString(data).Do(ctx)
	return err
}
