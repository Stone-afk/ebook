package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type TagElasticDAO struct {
	client *elastic.Client
}

func (dao *TagElasticDAO) InputBizTags(ctx context.Context, tag BizTags) error {
	docId := fmt.Sprintf("%d_%s_%d", tag.Uid, tag.Biz, tag.BizId)
	_, err := dao.client.Index().Index(TagIndexName).
		Id(docId).
		BodyJson(tag).Do(ctx)
	return err
}

func (dao *TagElasticDAO) Search(ctx context.Context, uid int64, biz string, keywords []string) ([]BizTags, error) {
	query := elastic.NewBoolQuery().Must(
		// 第一个条件，一定有 uid
		elastic.NewTermsQuery("uid", uid),
		// 第二个条件，biz 一定符合预期
		elastic.NewTermQuery("biz", biz),
		// 第三个条件，关键字命中了标签
		elastic.NewTermsQueryFromStrings("tags", keywords...))
	resp, err := dao.client.Search(TagIndexName).Query(query).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]BizTags, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var ele BizTags
		err = json.Unmarshal(hit.Source, &ele)
		if err != nil {
			return nil, err
		}
		res = append(res, ele)
	}
	return res, nil
}

func NewTagElasticDAO(client *elastic.Client) TagDAO {
	return &TagElasticDAO{client: client}
}
