package repository

import (
	"context"
	"database/sql"
	"ebook/cmd/comment/domain"
	"ebook/cmd/comment/repository/dao"
	"ebook/cmd/pkg/logger"
	"golang.org/x/sync/errgroup"
	"time"
)

type CachedCommentRepo struct {
	dao dao.CommentDAO
	l   logger.Logger
}

func (repo *CachedCommentRepo) toDomain(daoComment dao.Comment) domain.Comment {
	val := domain.Comment{
		Id: daoComment.Id,
		Commentator: domain.User{
			ID: daoComment.Uid,
		},
		Biz:     daoComment.Biz,
		BizID:   daoComment.BizID,
		Content: daoComment.Content,
		CTime:   time.UnixMilli(daoComment.Ctime),
		UTime:   time.UnixMilli(daoComment.Utime),
	}
	if daoComment.PID.Valid {
		val.ParentComment = &domain.Comment{
			Id: daoComment.PID.Int64,
		}
	}
	if daoComment.RootID.Valid {
		val.RootComment = &domain.Comment{
			Id: daoComment.RootID.Int64,
		}
	}
	return val
}

func (repo *CachedCommentRepo) toEntity(domainComment domain.Comment) dao.Comment {
	daoComment := dao.Comment{
		Id:      domainComment.Id,
		Uid:     domainComment.Commentator.ID,
		Biz:     domainComment.Biz,
		BizID:   domainComment.BizID,
		Content: domainComment.Content,
	}
	if domainComment.RootComment != nil {
		daoComment.RootID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.RootComment.Id,
		}
	}
	if domainComment.ParentComment != nil {
		daoComment.PID = sql.NullInt64{
			Valid: true,
			Int64: domainComment.ParentComment.Id,
		}
	}
	daoComment.Ctime = time.Now().UnixMilli()
	daoComment.Utime = time.Now().UnixMilli()
	return daoComment
}

func (repo *CachedCommentRepo) FindByBiz(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error) {
	daoComments, err := repo.dao.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Comment, 0, len(daoComments))
	// 只找三条
	var eg errgroup.Group
	downgraded := ctx.Value("downgraded") == "true"
	for _, d := range daoComments {
		d := d
		// 这两句不能放进去，因为并发操作 res 会有坑
		cm := repo.toDomain(d)
		res = append(res, cm)
		if downgraded {
			continue
		}
		eg.Go(func() error {
			// 只展示三条
			cm.Children = make([]domain.Comment, 0, 3)
			rs, er := repo.dao.FindRepliesByPid(ctx, cm.Id, 0, 3)
			if er != nil {
				// 我们认为这是一个可以容忍的错误
				repo.l.Error("查询子评论失败", logger.Error(err))
				return nil
			}
			for _, r := range rs {
				cm.Children = append(cm.Children, repo.toDomain(r))
			}
			return nil
		})
	}
	return res, eg.Wait()
}

func (repo *CachedCommentRepo) DeleteComment(ctx context.Context, comment domain.Comment) error {
	return repo.dao.Delete(ctx, dao.Comment{Id: comment.Id})
}

func (repo *CachedCommentRepo) CreateComment(ctx context.Context, comment domain.Comment) error {
	return repo.dao.Insert(ctx, repo.toEntity(comment))
}

func (repo *CachedCommentRepo) GetCommentByIds(ctx context.Context, ids []int64) ([]domain.Comment, error) {
	vals, err := repo.dao.FindOneByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	comments := make([]domain.Comment, 0, len(vals))
	for _, v := range vals {
		comment := repo.toDomain(v)
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repo *CachedCommentRepo) GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, error) {
	cs, err := repo.dao.FindRepliesByRid(ctx, rid, maxID, limit)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Comment, 0, len(cs))
	for _, cm := range cs {
		res = append(res, repo.toDomain(cm))
	}
	return res, nil
}

func NewCommentRepo(commentDAO dao.CommentDAO, l logger.Logger) CommentRepository {
	return &CachedCommentRepo{
		dao: commentDAO,
		l:   l,
	}
}
