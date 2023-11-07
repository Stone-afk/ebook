package cache

import (
	"context"
	"ebook/cmd/internal/domain"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	//go:embed lua/interative_incr_cnt.lua
	luaIncrCnt string
)

const (
	fieldReadCnt    = "read_cnt"
	fieldCollectCnt = "collect_cnt"
	fieldLikeCnt    = "like_cnt"
)

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/interactive.go -package=cachemocks -destination=/Users/stone/go_project/ebook/ebook/cmd/internal/repository/cache/mocks/interactive.mock.go
type InteractiveCache interface {
	// IncrReadCntIfPresent 如果在缓存中有对应的数据，就 +1
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error

	// Get 查询缓存中数据
	// 事实上，这里 liked 和 collected 是不需要缓存的
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
}

// 方案1
// key1 => map[string]int

// 方案2
// key1_read_cnt => 10
// key1_collect_cnt => 11
// key1_like_cnt => 13

type RedisInteractiveCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		client: client,
	}
}

func (r *RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	panic("")
}

func (r *RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	panic("")
}

func (r *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	panic("")
}

func (r *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	// 拿到的结果，可能自增成功了，可能不需要自增（key不存在）
	// 你要不要返回一个 error 表达 key 不存在？
	//res, err := r.client.Eval(ctx, luaIncrCnt,
	//	[]string{r.key(biz, bizId)},
	//	// read_cnt +1
	//	"read_cnt", 1).Int()
	//if err != nil {
	//	return err
	//}
	//if res == 0 {
	// 这边一般是缓存过期了
	//	return errors.New("缓存中 key 不存在")
	//}
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)}, fieldReadCnt, 1).Err()
}

func (r *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)}, fieldLikeCnt, 1).Err()
}

func (r *RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrCnt,
		[]string{r.key(biz, bizId)}, fieldLikeCnt, -1).Err()
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
