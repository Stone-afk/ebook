package ioc

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/tag/repository"
	"ebook/cmd/tag/repository/cache"
	"ebook/cmd/tag/repository/dao"
	"time"
)

func InitRepository(d dao.TagDAO, c cache.TagCache, l logger.Logger) repository.TagRepository {
	repo := repository.NewTagRepository(d, c, l)
	go func() {
		// 执行缓存预加载
		// 或者启动的环境变量
		// 启动参数控制
		// 或者借助配置中心的开关
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		// 也可以同步执行。但是在一些场景下，同步执行会占用很长的时间，所以可以考虑异步执行。
		err := repo.PreloadUserTags(ctx)
		if err != nil {
			l.Error("预加载缓存自定义标签列表失败", logger.Error(err))
		}
	}()
	return repo
}
