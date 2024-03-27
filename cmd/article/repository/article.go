package repository

import (
	"context"
	userv1 "ebook/cmd/api/proto/gen/user/v1"
	"ebook/cmd/article/domain"
	"ebook/cmd/article/repository/cache"
	"ebook/cmd/article/repository/dao"
	"ebook/cmd/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
)

type CachedArticleRepository struct {
	dao     dao.ArticleDAO
	cache   cache.ArticleCache
	userSvc userv1.UserServiceClient

	// SyncV1 用
	authorDAO dao.ArticleAuthorDAO
	readerDAO dao.ArticleReaderDAO

	// SyncV2 用
	db *gorm.DB
	l  logger.Logger
}

func NewCachedArticleRepository(dao dao.ArticleDAO,
	cache cache.ArticleCache,
	userSvc userv1.UserServiceClient,
	l logger.Logger) ArticleRepository {
	return &CachedArticleRepository{
		userSvc: userSvc,
		dao:     dao,
		cache:   cache,
		l:       l,
	}
}

func NewCachedArticleRepo(dao dao.ArticleDAO,
	cache cache.ArticleCache,
	userSvc userv1.UserServiceClient,
	l logger.Logger) *CachedArticleRepository {
	return &CachedArticleRepository{
		userSvc: userSvc,
		dao:     dao,
		cache:   cache,
		l:       l,
	}
}

func NewCachedArticleRepositoryV1(authorDAO dao.ArticleAuthorDAO,
	readerDAO dao.ArticleReaderDAO) ArticleRepository {
	return &CachedArticleRepository{
		authorDAO: authorDAO,
		readerDAO: readerDAO,
	}
}

func NewCachedArticleRepositoryV2(db *gorm.DB, l logger.Logger) ArticleRepository {
	return &CachedArticleRepository{
		db: db,
		l:  l,
	}
}

func (repo *CachedArticleRepository) Cache() cache.ArticleCache {
	return repo.cache
}

func (repo *CachedArticleRepository) ListPub(ctx context.Context, uTime time.Time, offset int, limit int) ([]domain.Article, error) {
	val, err := repo.dao.ListPubByUtime(ctx, uTime, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.PublishedArticle, domain.Article](val, func(idx int, src dao.PublishedArticle) domain.Article {
		// 偷懒写法
		return repo.ToDomain(dao.Article(src))
	}), nil
}

func (repo *CachedArticleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	res, err := repo.cache.GetPub(ctx, id)
	if err == nil {
		return res, err
	}
	// 读取线上库数据，如果你的 Content 被你放过去了 OSS 上，你就要让前端去读 Content 字段
	art, err := repo.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	resp, err := repo.userSvc.Profile(ctx, &userv1.ProfileRequest{
		Id: art.AuthorId,
	})
	if err != nil {
		return domain.Article{}, err
	}
	res = domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   resp.User.Id,
			Name: resp.User.Nickname,
		},
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

func (repo *CachedArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := repo.cache.Get(ctx, id)
	if err != nil {
		data, er := repo.dao.GetById(ctx, id)
		if er != nil {
			return domain.Article{}, err
		}
		repo.l.Error("查询缓存文章失败", logger.Int64("id", id), logger.Error(err))
		return repo.ToDomain(data), nil
	}
	return art, nil
}

func (repo *CachedArticleRepository) List(ctx context.Context, authorId int64, offset, limit int) ([]domain.Article, error) {
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
	res := slice.Map[dao.Article, domain.Article](arts,
		func(idx int, src dao.Article) domain.Article {
			return repo.ToDomain(src)
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

func (repo *CachedArticleRepository) preCache(ctx context.Context, arts []domain.Article) {
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

func (repo *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 直接 defer Rollback
	// 如果我们后续 Commit 了，这里会得到一个错误，但是没关系
	defer tx.Rollback()
	authorDAO := dao.NewGORMArticleDAO(tx)
	readerDAO := dao.NewGORMArticleReaderDAO(tx)

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
	err = readerDAO.UpsertV2(ctx, dao.PublishedArticle(artn))
	if err != nil {
		// 依赖于 defer 来 rollback
		return 0, err
	}
	tx.Commit()
	return artn.Id, nil
}

func (repo *CachedArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
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

func (repo *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := repo.dao.Sync(ctx, repo.toEntity(art))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *CachedArticleRepository) SyncStatus(ctx context.Context, uid, id int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, uid, id, status.ToUint8())
}

func (repo *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(art))
}

func (repo *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(art))
}

func (repo *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}

func (repo *CachedArticleRepository) ToDomain(art dao.Article) domain.Article {
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

type GrpcAuthorRepository struct {
	userService userv1.UserServiceClient
	dao         dao.ArticleDAO
}

func NewGrpcAuthorRepository(articleDao dao.ArticleDAO, userService userv1.UserServiceClient) AuthorRepository {
	return &GrpcAuthorRepository{
		userService: userService,
		dao:         articleDao,
	}
}

func (g *GrpcAuthorRepository) FindAuthor(ctx context.Context, id int64) (domain.Author, error) {
	art, err := g.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Author{}, nil
	}
	u, err := g.userService.Profile(ctx, &userv1.ProfileRequest{
		Id: art.AuthorId,
	})
	if err != nil {
		return domain.Author{}, err
	}
	return domain.Author{
		Id:   u.User.Id,
		Name: u.User.Nickname,
	}, nil
}
