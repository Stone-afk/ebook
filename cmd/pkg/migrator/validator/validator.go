package validator

import (
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/migrator"
	"gorm.io/gorm"
)

// Validator T 必须实现了 Entity 接口
type Validator[T migrator.Entity] struct {
	// 校验，以 XXX 为准，
	base *gorm.DB
	// 校验的是谁的数据
	target    *gorm.DB
	l         logger.Logger
	direction string

	batchSize int
}
