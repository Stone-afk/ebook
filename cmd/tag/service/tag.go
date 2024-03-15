package service

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/tag/domain"
	"ebook/cmd/tag/events"
	"ebook/cmd/tag/repository"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type tagService struct {
	repo     repository.TagRepository
	l        logger.Logger
	producer events.Producer
}

func (s *tagService) CreateTag(ctx context.Context, uid int64, name string) (int64, error) {
	return s.repo.CreateTag(ctx, domain.Tag{
		Uid:  uid,
		Name: name,
	})
}

func (s *tagService) AttachTags(ctx context.Context, uid int64, biz string, bizId int64, tags []int64) error {
	err := s.repo.BindTagToBiz(ctx, uid, biz, bizId, tags)
	if err != nil {
		return err
	}
	// 异步发送
	go func() {
		ts, er := s.repo.GetTagsById(ctx, tags)
		if er != nil {
			// 记录日志
			s.l.Error("获取自定义标签失败", logger.Error(er))
		}
		// 这里要根据 tag_index 的结构来定义
		// 同样要注意顺序，即同一个用户对同一个资源打标签的顺序，
		// 是不能乱的
		pctx, cancel := context.WithTimeout(context.Background(), time.Second)
		er = s.producer.ProduceSyncEvent(pctx, events.BizTags{
			Uid:   uid,
			Biz:   biz,
			BizId: bizId,
			Tags: slice.Map(ts, func(idx int, src domain.Tag) string {
				return src.Name
			}),
		})
		cancel()
		if er != nil {
			// 记录日志
			s.l.Error("同步 kafka 事件消息失败", logger.Error(er))
		}
	}()
	return err
}

func (s *tagService) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	return s.repo.GetTags(ctx, uid)
}

func (s *tagService) GetBizTags(ctx context.Context, uid int64, biz string, bizId int64) ([]domain.Tag, error) {
	return s.repo.GetBizTags(ctx, uid, biz, bizId)
}

func NewTagService(repo repository.TagRepository, l logger.Logger) TagService {
	return &tagService{
		repo: repo,
		l:    l,
	}
}
