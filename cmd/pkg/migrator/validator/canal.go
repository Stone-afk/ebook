package validator

import "ebook/cmd/pkg/migrator"

type CanalIncrValidator[T migrator.Entity] struct {
	baseValidator
}
