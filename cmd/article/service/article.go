package service

import (
	"context"
	"ebook/cmd/article/domain"
	"ebook/cmd/article/events"
	"ebook/cmd/article/repository"
	"ebook/cmd/pkg/logger"
	"fmt"
	"strconv"
	"time"
)

type articleService struct {
	// 1. 在 service 这一层使用两个 repository
	authorRepo repository.ArticleAuthorRepository
	readerRepo repository.ArticleReaderRepository

	// 2. 在 repo 里面处理制作库和线上库
	// 1 和 2 是互斥的，不会同时存在
	repo                    repository.ArticleRepository
	readEventProducer       events.ReadEventProducer
	syncSearchEventProducer events.SyncSearchEventProducer
	feedEventProducer       events.FeedEventProducer
	log                     logger.Logger

	ch chan readInfo
}

func NewArticleService(repo repository.ArticleRepository,
	readEventProducer events.ReadEventProducer,
	syncSearchEventProducer events.SyncSearchEventProducer,
	feedEventProducer events.FeedEventProducer,
	l logger.Logger) ArticleService {
	return &articleService{
		repo:                    repo,
		log:                     l,
		readEventProducer:       readEventProducer,
		syncSearchEventProducer: syncSearchEventProducer,
		feedEventProducer:       feedEventProducer,
	}
}

//func NewArticleServiceV2(repo repository.ArticleRepository,
//	l logger.Logger,
//	producer events.Producer) ArticleService {
//	ch := make(chan readInfo, 10)
//	go func() {
//		for {
//			uids := make([]int64, 0, 10)
//			aids := make([]int64, 0, 10)
//			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//			for i := 0; i < 10; i++ {
//				select {
//				case info, ok := <-ch:
//					if !ok {
//						cancel()
//						return
//					}
//					uids = append(uids, info.uid)
//					aids = append(aids, info.aid)
//				case <-ctx.Done():
//					break
//				}
//			}
//			cancel()
//			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
//			err := producer.ProduceReadEventV1(ctx, events.ReadEventV1{
//				Uids: uids,
//				Aids: aids,
//			})
//			if err == nil {
//				l.Error("发送读者阅读事件失败")
//			}
//			cancel()
//		}
//	}()
//	return &articleService{
//		repo:     repo,
//		producer: producer,
//		log:      l,
//		ch:       ch,
//	}
//}

type readInfo struct {
	uid int64
	aid int64
}

func NewArticleServiceV1(
	authorRepo repository.ArticleAuthorRepository,
	readerRepo repository.ArticleReaderRepository,
	producer events.ReadEventProducer,
	l logger.Logger) ArticleService {
	return &articleService{
		authorRepo:        authorRepo,
		readerRepo:        readerRepo,
		log:               l,
		readEventProducer: producer,
	}
}

func (svc *articleService) ListPub(ctx context.Context,
	uTime time.Time, offset, limit int) ([]domain.Article, error) {
	return svc.repo.ListPub(ctx, uTime, offset, limit)
}

func (svc *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return svc.repo.GetById(ctx, id)
}

func (svc *articleService) GetPublishedById(ctx context.Context, id, userId int64) (domain.Article, error) {
	art, err := svc.repo.GetPublishedById(ctx, id)
	if err == nil {
		go func() {
			// 生产者也可以通过改批量来提高性能
			er := svc.readEventProducer.ProduceReadEvent(
				ctx,
				events.ReadEvent{
					// 即便你的消费者要用 art 的里面的数据，
					// 让它去查询，你不要在 event 里面带
					Uid: userId,
					Aid: id,
				})
			if er != nil {
				svc.log.Error("发送读者阅读事件失败", logger.Error(er))
			}
		}()
	}
	return art, err
}

func (svc *articleService) List(ctx context.Context, authorId int64, offset, limit int) ([]domain.Article, error) {
	return svc.repo.List(ctx, authorId, offset, limit)
}

func (svc *articleService) Withdraw(ctx context.Context, uid, id int64) error {
	// art.Status = domain.ArticleStatusPrivate 然后直接把整个 art 往下传
	err := svc.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusPrivate)
	if err == nil {
		go func() {
			art, err := svc.repo.GetById(ctx, id)
			if err != nil {
				svc.log.Error(fmt.Sprintf("发送同步搜索文章事件失败: 查询文章出错 %s", err.Error()))
			} else {
				svc.sendSyncEventToSearch(ctx, art)
			}
		}()
	}
	return err
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := svc.repo.Update(ctx, art)
		return art.Id, err
	}
	artId, err := svc.repo.Create(ctx, art)
	if err == nil {
		go func() {
			svc.sendSyncEventToSearch(ctx, art)
		}()
	}
	return artId, err
}

func (svc *articleService) Publish(ctx context.Context,
	art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	artId, err := svc.repo.Sync(ctx, art)
	if err == nil {
		go func() {
			svc.sendSyncEventToSearch(ctx, art)
			svc.sendFeedEvent(ctx, art)
		}()
	}
	return artId, err
}

func (svc *articleService) sendSyncEventToSearch(ctx context.Context, art domain.Article) {
	evt := events.ArticleEvent{
		Id:      art.Id,
		Title:   art.Title,
		Status:  int32(art.Status),
		Content: art.Content,
	}
	er := svc.syncSearchEventProducer.ProduceSyncEvent(ctx, evt)
	if er != nil {
		svc.log.Error("ProduceSyncEvent 发送同步搜索文章事件失败", logger.Error(er))
		er = svc.syncSearchEventProducer.ProduceStandardSyncEvent(ctx, evt)
		if er != nil {
			svc.log.Error("ProduceStandardSyncEvent 发送同步搜索文章事件失败", logger.Error(er))
		}
	}
}

func (svc *articleService) sendFeedEvent(ctx context.Context, art domain.Article) {
	evt := events.FeedEvent{
		Type: "article_event",
		Metadata: map[string]string{
			"uid": strconv.FormatInt(art.Author.Id, 10),
			"aid": strconv.FormatInt(art.Id, 10),
		},
	}
	er := svc.feedEventProducer.ProduceFeedEvent(ctx, evt)
	if er != nil {
		svc.log.Error("ProduceFeedEvent 发送feed流事件失败", logger.Error(er))
		er = svc.feedEventProducer.ProduceStandardFeedEvent(ctx, evt)
		if er != nil {
			svc.log.Error("ProduceStandardFeedEvent 发送feed流事件失败", logger.Error(er))
		}
	}
}

// PublishV1 基于使用两种 repository 的写法
func (svc *articleService) PublishV1(ctx context.Context,
	art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	// 这一段逻辑其实就是 Save
	if art.Id == 0 {
		id, err = svc.authorRepo.Create(ctx, art)
	} else {
		err = svc.authorRepo.Update(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	// 保持制作库和线上库的 ID 是一样的。
	art.Id = id
	for i := 0; i < 3; i++ {
		err = svc.readerRepo.Save(ctx, art)
		if err == nil {
			break
		}
		svc.log.Error("部分失败：保存数据到线上库失败",
			logger.Int64("art_id", id),
			logger.Error(err))
	}
	// 在接入了 metrics 或者 tracing 之后，
	// 这边要进一步记录必要的DEBUG信息。
	if err != nil {
		svc.log.Error("全部失败：保存数据到线上库重试都失败了",
			logger.Int64("art_id", id),
			logger.Error(err))
	}
	return id, nil
}
