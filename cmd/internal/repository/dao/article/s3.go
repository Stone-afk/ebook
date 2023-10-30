package article

import (
	"context"
	"ebook/cmd/internal/domain"
	_ "github.com/aws/aws-sdk-go"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"gorm.io/gorm"
)

var statusPrivate = domain.ArticleStatusPrivate.ToUint8()

type S3DAO struct {
	oss *s3.S3
	// 通过组合 GORMArticleDAO 来简化操作
	// 当然在实践中，你是不太会有组合的机会
	GORMArticleDAO
	bucket *string
}

// NewOssDAO 因为组合 GORMArticleDAO 是一个内部实现细节
// 所以这里要直接传入 DB
func NewOssDAO(oss *s3.S3, db *gorm.DB) ArticleDAO {
	return &S3DAO{
		oss: oss,
		// 你也可以考虑利用依赖注入来传入。
		// 但是事实上这个很少变，所以你可以延迟到必要的时候再注入
		bucket: ekit.ToPtr[string]("ebook-1314583317"),
		GORMArticleDAO: GORMArticleDAO{
			db: db,
		},
	}
}

func (dao *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	panic("implement me")
}

func (dao *S3DAO) SyncStatus(ctx context.Context, authorId, id int64, status uint8) error {
	panic("implement me")
}

func (dao *S3DAO) Insert(ctx context.Context, art Article) (int64, error) {
	panic("implement me")
}

func (dao *S3DAO) UpdateById(ctx context.Context, art Article) error {
	panic("implement me")
}
