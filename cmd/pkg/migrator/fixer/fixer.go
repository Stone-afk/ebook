package fixer

import (
	"context"
	"ebook/cmd/pkg/migrator"
	"ebook/cmd/pkg/migrator/events"
	"gorm.io/gorm"
)

type Fixer[T migrator.Entity] struct {
	base    *gorm.DB
	target  *gorm.DB
	columns []string
}

// Fix 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	panic("")
}

// FixV1 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) FixV1(ctx context.Context, evt events.InconsistentEvent) error {
	panic("")
}

// FixV2 最一了百了的写法
// 不管三七二十一，我TM直接覆盖
// 把 event 当成一个触发器，不依赖的 event 的具体内容（ID 必须不可变）
// 修复这里，也改成批量？？
func (f *Fixer[T]) FixV2(ctx context.Context, evt events.InconsistentEvent) error {
	panic("")
}
