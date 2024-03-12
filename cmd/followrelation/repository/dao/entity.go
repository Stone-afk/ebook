package dao

const (
	FollowRelationStatusUnknown uint8 = iota
	FollowRelationStatusActive
	FollowRelationStatusInactive
)

type FollowRelation struct {
	ID int64 `gorm:"primaryKey,autoIncrement,column:id"`
	// 关注人id
	Follower int64 `gorm:"type:int(11);not null;uniqueIndex:follower_followee"`
	// 被关注人id
	Followee int64 `gorm:"type:int(11);not null;uniqueIndex:follower_followee"`

	Status uint8

	// 这里你可以根据自己的业务来增加字段，比如说
	// 关系类型，可以搞些什么普通关注，特殊关注
	// Type int64 `gorm:"column:type;type:int(11);comment:关注类型 0-普通关注"`
	// 备注
	// Note string `gorm:"column:remark;type:varchar(255);"`
	// 创建时间
	Ctime int64
	Utime int64
}

// UserRelation 另外一种设计方案，但是不建议这么做
type UserRelation struct {
	ID     int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid1   int64 `gorm:"column:uid1;type:int(11);not null;uniqueIndex:user_contact_index"`
	Uid2   int64 `gorm:"column:uid2;type:int(11);not null;uniqueIndex:user_contact_index"`
	Block  bool  // 拉黑
	Mute   bool  // 屏蔽
	Follow bool  // 关注
}

type FollowStatics struct {
	ID  int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid int64 `gorm:"unique"`
	// 有多少粉丝
	Followers int64
	// 关注了多少人
	Followees int64

	Utime int64
	Ctime int64
}
