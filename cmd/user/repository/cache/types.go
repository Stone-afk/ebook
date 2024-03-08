package cache

import (
	"context"
	"ebook/cmd/user/domain"
	"github.com/redis/go-redis/v9"
)

// ErrKeyNotExist 因为我们目前还是只有一个实现，所以可以保持用别名
var ErrKeyNotExist = redis.Nil

//go:generate mockgen -source=/Users/stone/go_project/ebook/ebook/cmd/user/repository/cache/types.go -package=cachemocks -destination=/Users/stone/go_project/ebook/ebook/cmd/user/repository/cache/mocks/user.mock.go
type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	Delete(ctx context.Context, id int64) error
}
