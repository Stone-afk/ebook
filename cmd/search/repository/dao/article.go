package dao

import (
	"context"
	"ebook/cmd/search/repository/dao/es"
	"github.com/ecodeclub/ekit/slice"
	"github.com/olivere/elastic/v7"
	"strconv"
	"strings"
)

type ArticleElasticDAO struct {
	client *elastic.Client
}

func (dao *ArticleElasticDAO) InputArticle(ctx context.Context, article Article) error {
	_, err := dao.client.Index().Index(ArticleIndexName).
		// 为什么要指定 ID？
		// 确保后面文章更新的时候，我们这里产生类似的两条数据，而是更新了数据
		Id(strconv.FormatInt(article.Id, 10)).
		BodyJson(article).Do(ctx)
	return err
}

func (dao *ArticleElasticDAO) Search(ctx context.Context, tagArtIds []int64, keywords []string) ([]Article, error) {
	queryString := strings.Join(keywords, " ")
	// 文章，标题或者内容任何一个匹配上
	// 并且状态 status 必须是已发表的状态
	// status 精确查找
	statusTerm := elastic.NewTermQuery("status", 2)
	// 标签命中
	tagArtIdAnys := slice.Map(tagArtIds, func(idx int, src int64) any {
		return src
	})
	// 内容或者标题，模糊查找（match）
	title := elastic.NewMatchQuery("title", queryString)
	content := elastic.NewMatchQuery("content", queryString)
	or := elastic.NewBoolQuery().Should(title, content)
	if len(tagArtIds) > 0 {
		tag := elastic.NewTermsQuery("id", tagArtIdAnys...).
			Boost(2.0)
		or = or.Should(tag)
	}
	and := elastic.NewBoolQuery().Must(statusTerm, or)
	return es.NewSearcher[Article](dao.client, ArticleIndexName).
		Query(and).Do(ctx)

	//resp, err := dao.client.Search(ArticleIndexName).Query(and).Do(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//var res []Article
	//for _, hit := range resp.Hits.Hits {
	//	var art Article
	//	err = json.Unmarshal(hit.Source, &art)
	//	if err != nil {
	//		return nil, err
	//	}
	//	res = append(res, art)
	//}
	//return res, nil
}

func NewArticleElasticDAO(client *elastic.Client) ArticleDAO {
	return &ArticleElasticDAO{client: client}
}
