package service

import (
	"context"
	articlev1 "ebook/cmd/api/proto/gen/article/v1"
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	"ebook/cmd/ranking/domain"
	"ebook/cmd/ranking/repository"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"time"
)

// BatchRankingService 分批计算
type BatchRankingService struct {
	// 为了测试，不得已暴露出去
	BatchSize int
	N         int

	//artSvc  ArticleService
	artSvc  articlev1.ArticleServiceClient
	intrSvc intrv1.InteractiveServiceClient
	// 将来扩展，以及支持测试
	repo repository.RankingRepository
	// scoreFunc 不能返回负数
	scoreFunc func(likeCnt int64, utime time.Time) float64
}

func NewBatchRankingService(
	intrSvc intrv1.InteractiveServiceClient,
	artSvc articlev1.ArticleServiceClient,
	repo repository.RankingRepository) RankingService {
	res := &BatchRankingService{
		intrSvc:   intrSvc,
		artSvc:    artSvc,
		repo:      repo,
		BatchSize: 100,
		N:         100,
	}
	res.scoreFunc = res.score
	return res
}

func (svc *BatchRankingService) RankTopN(ctx context.Context) error {
	arts, err := svc.rankTopN(ctx)
	if err != nil {
		return err
	}
	// 准备放到缓存里面
	return svc.repo.ReplaceTopN(ctx, arts)
}

func (svc *BatchRankingService) rankTopN(ctx context.Context) ([]domain.Article, error) {
	now := time.Now()
	// 只计算七天内的，因为超过七天的我们可以认为绝对不可能成为热榜了
	// 如果一个批次里面 utime 最小已经是七天之前的，我们就中断当前计算
	ddl := now.Add(-time.Hour * 24 * 7)
	// 先拿一批数据
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	// 这里可以用非并发安全
	topN := queue.NewConcurrentPriorityQueue[Score](svc.N,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			} else {
				return -1
			}
		})
	for {
		// 这里拿了一批
		arts, err := svc.artSvc.ListPub(ctx, &articlev1.ListPubRequest{
			StartTime: timestamppb.New(now),
			Offset:    int32(offset),
			Limit:     int32(svc.BatchSize),
		})
		if err != nil {
			return nil, err
		}
		// 转化成 domain Article
		domainArts := make([]domain.Article, 0, len(arts.Articles))
		for _, art := range arts.Articles {
			domainArts = append(domainArts, articleToDomain(art))
		}

		ids := slice.Map[domain.Article, int64](domainArts,
			func(idx int, src domain.Article) int64 {
				return src.Id
			})
		// 要去找到对应的点赞数据
		resp, err := svc.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
			Biz: "article", Ids: ids,
		})
		if err != nil {
			return nil, err
		}
		intrs := resp.Intrs
		if len(intrs) == 0 {
			return nil, errors.New("没有数据")
		}
		// 合并计算 score
		// 排序
		for _, art := range domainArts {
			intr := intrs[art.Id]
			//if !ok {
			//	// 你都没有，肯定不可能是热榜
			//	continue
			//}
			score := svc.scoreFunc(intr.LikeCnt, art.Utime)
			// 要考虑，这个 score 在不在前一百名
			// 拿到热度最低的
			err = topN.Enqueue(Score{
				art:   art,
				score: score,
			})
			// 这种写法，要求 topN 已经满了
			if err == queue.ErrOutOfCapacity {
				val, _ := topN.Dequeue()
				if val.score < score {
					err = topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				} else {
					_ = topN.Enqueue(val)
				}
			}
		}
		// 一批已经处理完了，问题来了，我要不要进入下一批？我怎么知道还有没有？
		if len(domainArts) == 0 || len(domainArts) < svc.BatchSize ||
			domainArts[len(domainArts)-1].Utime.Before(ddl) {
			// 我这一批都没取够，我当然可以肯定没有下一批了
			break
		}
		// 这边要更新 offset
		offset = offset + len(domainArts)
	}
	// 最后得出结果
	res := make([]domain.Article, svc.N)
	for i := svc.N - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			// 说明取完了，不够 n
			break
		}
		res[i] = val.art
	}
	return res, nil
}

func (svc *BatchRankingService) TopN(ctx context.Context) ([]domain.Article, error) {
	return svc.repo.GetTopN(ctx)
}

// 这里不需要提前抽象算法，因为正常一家公司的算法都是固定的，不会今天切换到这里，明天切换到那里
func (svc *BatchRankingService) score(likeCnt int64, utime time.Time) float64 {
	// 这个 factor 也可以做成一个参数
	const factor = 1.5
	return float64(likeCnt-1) /
		math.Pow(time.Since(utime).Hours()+2, factor)
}

func articleToDomain(article *articlev1.Article) domain.Article {
	domainArticle := domain.Article{}
	if article != nil {
		domainArticle.Id = article.GetId()
		domainArticle.Title = article.GetTitle()
		domainArticle.Status = domain.ArticleStatus(article.Status)
		domainArticle.Content = article.Content
		domainArticle.Author = domain.Author{
			Id:   article.GetAuthor().GetId(),
			Name: article.GetAuthor().GetName(),
		}
		domainArticle.Ctime = article.Ctime.AsTime()
		domainArticle.Utime = article.Utime.AsTime()
	}
	return domainArticle
}
