package repository

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/repository/cache"
	"ebook/cmd/internal/repository/dao/article"
	"ebook/cmd/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

// repository 还是要用来操作缓存和DAO
// 事务概念应该在 DAO 这一层

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/article.go -package=repomocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/mocks/article.mock.go

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	// Sync 本身要求先保存到制作库，再同步到线上库
	Sync(ctx context.Context, art domain.Article) (int64, error)
	// SyncStatus 仅仅同步状态
	SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error
	List(ctx context.Context, author int64, offset, limit int) ([]domain.Article, error)
	ListPub(ctx context.Context, uTime time.Time, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type articleRepository struct {
	dao      article.ArticleDAO
	cache    cache.ArticleCache
	userRepo UserRepository

	// SyncV1 用
	authorDAO article.ArticleAuthorDAO
	readerDAO article.ArticleReaderDAO

	// SyncV2 用
	db *gorm.DB
	l  logger.Logger
}

func NewArticleRepository(dao article.ArticleDAO,
	cache cache.ArticleCache,
	userRepo UserRepository,
	l logger.Logger) ArticleRepository {
	return &articleRepository{
		userRepo: userRepo,
		dao:      dao,
		cache:    cache,
		l:        l,
	}
}

func NewArticleRepositoryV1(authorDAO article.ArticleAuthorDAO,
	readerDAO article.ArticleReaderDAO) ArticleRepository {
	return &articleRepository{
		authorDAO: authorDAO,
		readerDAO: readerDAO,
	}
}

func NewArticleRepositoryV2(db *gorm.DB, l logger.Logger) ArticleRepository {
	return &articleRepository{
		db: db,
		l:  l,
	}
}

func (repo *articleRepository) ListPub(ctx context.Context, uTime time.Time, offset int, limit int) ([]domain.Article, error) {
	val, err := repo.dao.ListPubByUtime(ctx, uTime, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[article.PublishedArticle, domain.Article](val, func(idx int, src article.PublishedArticle) domain.Article {
		// 偷懒写法
		return repo.toDomain(article.Article(src))
	}), nil
}

func (repo *articleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	// 读取线上库数据，如果你的 Content 被你放过去了 OSS 上，你就要让前端去读 Content 字段
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	// 在这边要组装 user 了，适合单体应用
	usr, err := repo.userRepo.FindById(ctx, art.AuthorId)
	if err != nil {
		return domain.Article{}, err
	}
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.Nickname,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	// 也可以同步
	go func() {
		if err = repo.cache.SetPub(ctx, res); err != nil {
			repo.l.Error("缓存已发表文章失败",
				logger.Error(err), logger.Int64("aid", art.Id))
		}
	}()
	return res, nil
}

func (repo *articleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := repo.cache.Get(ctx, id)
	if err != nil {
		data, er := repo.dao.GetById(ctx, id)
		if er != nil {
			return domain.Article{}, err
		}
		repo.l.Error("查询缓存文章失败", logger.Int64("id", id), logger.Error(err))
		return repo.toDomain(data), nil
	}
	return art, nil
}

func (repo *articleRepository) List(ctx context.Context, authorId int64, offset, limit int) ([]domain.Article, error) {
	// 只有第一页才走缓存，并且假定一页只有 100 条
	// 也就是说，如果前端允许创作者调整页的大小
	// 那么只有 100 这个页大小这个默认情况下，会走索引
	if offset == 0 && limit == 100 {
		data, err := repo.cache.GetFirstPage(ctx, authorId)
		if err == nil {
			// 提前准备文章内容详情的缓存  一般都是让调用者来控制是否异步。
			go func() {
				repo.preCache(ctx, data)
			}()
			return data, nil
		}
		if err != cache.ErrKeyNotExist {
			repo.l.Error("查询缓存文章失败",
				logger.Int64("author", authorId), logger.Error(err))
		}
	}
	arts, err := repo.dao.GetByAuthor(ctx, authorId, offset, limit)
	if err != nil {
		return nil, err
	}
	res := slice.Map[article.Article, domain.Article](arts,
		func(idx int, src article.Article) domain.Article {
			return repo.toDomain(src)
		})
	// 提前准备文章内容详情的缓存  一般都是让调用者来控制是否异步。
	go func() {
		repo.preCache(ctx, res)
	}()
	// 这个也可以做成异步的
	err = repo.cache.SetFirstPage(ctx, authorId, res)
	if err != nil {
		repo.l.Error("刷新第一页文章的缓存失败",
			logger.Int64("author", authorId), logger.Error(err))
	}
	return res, nil
}

func (repo *articleRepository) preCache(ctx context.Context, arts []domain.Article) {
	// 1MB
	const contentSizeThreshold = 1024 * 1024
	// 只缓存第一篇文章
	if len(arts) > 0 && len(arts[0].Content) <= contentSizeThreshold {
		// 你也可以记录日志
		if err := repo.cache.Set(ctx, arts[0]); err != nil {
			repo.l.Error("提前准备缓存失败", logger.Error(err))
		}
	}
}

func (repo *articleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 直接 defer Rollback
	// 如果我们后续 Commit 了，这里会得到一个错误，但是没关系
	defer tx.Rollback()
	authorDAO := article.NewGORMArticleDAO(tx)
	readerDAO := article.NewGORMArticleReaderDAO(tx)

	// 下面代码和 SyncV1 一模一样
	artn := repo.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = authorDAO.Insert(ctx, artn)
		if err != nil {
			return 0, err
		}
	} else {
		err = authorDAO.UpdateById(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = readerDAO.UpsertV2(ctx, article.PublishedArticle(artn))
	if err != nil {
		// 依赖于 defer 来 rollback
		return 0, err
	}
	tx.Commit()
	return artn.Id, nil
}

func (repo *articleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	artn := repo.toEntity(art)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = repo.authorDAO.Create(ctx, artn)
		if err != nil {
			return 0, err
		}
	} else {
		err = repo.authorDAO.UpdateById(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	artn.Id = id
	err = repo.readerDAO.Upsert(ctx, artn)
	return id, err
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *articleRepository) SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, uid, id, status.ToUint8())
}

func (repo *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(art))
}

func (repo *articleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}

func (repo *articleRepository) toDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Status: domain.ArticleStatus(art.Status),
	}
}
