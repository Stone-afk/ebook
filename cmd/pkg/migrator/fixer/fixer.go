package fixer

import (
	"ebook/cmd/pkg/migrator"
	"gorm.io/gorm"
)

type Fixer[T migrator.Entity] struct {
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}
