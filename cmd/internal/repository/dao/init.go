package dao

import (
	"ebook/cmd/internal/repository/dao/user"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&user.User{})
}
