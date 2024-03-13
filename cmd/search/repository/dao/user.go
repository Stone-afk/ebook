package dao

import (
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

type UserElasticDAO struct {
	client *elastic.Client
}

func (dao *UserElasticDAO) InputUser(ctx context.Context, user User) error {
	_, err := dao.client.Index().
		Index(UserIndexName).
		Id(strconv.FormatInt(user.Id, 10)).
		BodyJson(user).Do(ctx)
	return err
}

func (dao *UserElasticDAO) Search(ctx context.Context, keywords []string) ([]User, error) {
	// 纯粹是因为前面已经预处理了输入
	queryString := strings.Join(keywords, " ")
	// 昵称命中就可以的
	resp, err := dao.client.Search(UserIndexName).
		Query(elastic.NewMatchQuery("nickname", queryString)).Do(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]User, 0, resp.Hits.TotalHits.Value)
	for _, hit := range resp.Hits.Hits {
		var u User
		err = json.Unmarshal(hit.Source, &u)
		if err != nil {
			return nil, err
		}
		res = append(res, u)
	}
	return res, nil
}

func NewUserElasticDAO(client *elastic.Client) UserDAO {
	return &UserElasticDAO{
		client: client,
	}
}
