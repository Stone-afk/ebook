package repository

import (
	"context"
	"ebook/cmd/search/domain"
	"ebook/cmd/search/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type tagRepository struct {
	dao dao.TagDAO
}

func (repo *tagRepository) InputTag(ctx context.Context, msg domain.BizTags) error {
	return repo.dao.InputTag(ctx, dao.BizTags{
		Uid:   msg.Uid,
		Biz:   msg.Biz,
		BizId: msg.BizId,
		Tags:  msg.Tags,
	})
}

func (repo *tagRepository) SearchTag(ctx context.Context, uid int64, biz string, keywords []string) ([]domain.BizTags, error) {
	tags, err := repo.dao.Search(ctx, uid, biz, keywords)
	if err != nil {
		return nil, err
	}
	return slice.Map(tags, func(idx int, src dao.BizTags) domain.BizTags {
		return domain.BizTags{
			Uid:   src.Uid,
			Biz:   src.Biz,
			BizId: src.BizId,
			Tags:  src.Tags,
		}
	}), nil
}

func NewTagRepository(d dao.TagDAO) TagRepository {
	return &tagRepository{
		dao: d,
	}
}
