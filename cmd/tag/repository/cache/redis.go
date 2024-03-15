package cache

import (
	"context"
	"ebook/cmd/tag/domain"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type RedisTagCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *RedisTagCache) GetTags(ctx context.Context, uid int64) ([]domain.Tag, error) {
	key := r.userTagsKey(uid)
	//data, err := r.client.LRange(ctx, key, 0, -1).Result()
	//if err != nil {
	//	return nil, err
	//}
	//res := make([]domain.Tag, 0, len(data))
	//for _, ele := range data {
	//	var t domain.Tag
	//	err = json.Unmarshal([]byte(ele), &t)
	//	if err != nil {
	//		return nil, err
	//	}
	//	res = append(res, t)
	//}

	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	res := make([]domain.Tag, 0, len(data))
	for _, ele := range data {
		var t domain.Tag
		err = json.Unmarshal([]byte(ele), &t)
		if err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, nil
}

func (r *RedisTagCache) Append(ctx context.Context, uid int64, tags ...domain.Tag) error {
	// data := make([]any, 0, len(tags))
	key := r.userTagsKey(uid)
	pip := r.client.Pipeline()
	for _, tag := range tags {
		val, err := json.Marshal(tag)
		if err != nil {
			return err
		}
		// data = append(data, val)
		pip.HMSet(ctx, key, strconv.FormatInt(tag.Id, 10), val)
	}
	// 利用 pipeline 来执行，性能好一点

	// pip.RPush(ctx, key, data...) 使用 List 在多个实例启动时预加载数据可能导致冲突，所以改成 hash
	// 无法辨别 key 是不是已经有过期时间，
	// 如果这个 key 的第一个 tag，是没有过期时间
	// 可以不设置过期时间
	pip.Expire(ctx, key, r.expiration)
	_, err := pip.Exec(ctx)
	return err
}

func (r *RedisTagCache) DelTags(ctx context.Context, uid int64) error {
	return r.client.Del(ctx, r.userTagsKey(uid)).Err()
}

func (r *RedisTagCache) userTagsKey(uid int64) string {
	return fmt.Sprintf("tag:user_tags:%d", uid)
}

func NewRedisTagCache(client redis.Cmdable) TagCache {
	return &RedisTagCache{
		client: client,
	}
}
