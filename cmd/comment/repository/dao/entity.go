package dao

import "database/sql"

type Comment struct {
	Id int64 `gorm:"column:id;primaryKey" json:"id"`
	// 发表评论的用户
	Uid int64 `gorm:"column:uid;index" json:"uid"`
	// 发表评论的业务类型
	Biz string `gorm:"column:biz;index:biz_type_id" json:"biz"`
	// 对应的业务ID
	BizID int64 `gorm:"column:biz_id;index:biz_type_id" json:"bizID"`
	// 根评论为0表示一级评论
	RootID sql.NullInt64 `gorm:"column:root_id;index" json:"rootID"`
	// 父级评论
	PID sql.NullInt64 `gorm:"column:pid;index" json:"pid"`
	// 外键 用于级联删除
	ParentComment *Comment `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"`
	// 评论内容
	Content string `gorm:"type:text;column:content" json:"content"`
	// 创建时间
	Ctime int64 `gorm:"column:ctime;" json:"ctime"`
	// 更新时间
	Utime int64 `gorm:"column:utime;" json:"utime"`
}

func (*Comment) TableName() string {
	return "comments"
}
