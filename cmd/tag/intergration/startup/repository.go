package startup

import (
	"ebook/cmd/pkg/logger"
	"ebook/cmd/tag/repository"
	"ebook/cmd/tag/repository/cache"
	"ebook/cmd/tag/repository/dao"
)

func InitRepository(d dao.TagDAO, c cache.TagCache, l logger.Logger) repository.TagRepository {
	return repository.NewTagRepository(d, c, l)
}
