package dao

type FeedPullEvent struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	UID     int64  `gorm:"column:uid;type:int(11);not null;"`
	Type    string `gorm:"column:type;type:varchar(255);comment:类型"`
	Content string `gorm:"column:content;type:text;"`
	// 发生时间
	Ctime int64 `gorm:"column:ctime;comment:发生时间"`
}

type FeedPushEvent struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	UID     int64  `gorm:"column:uid;type:int(11);not null;"`
	Type    string `gorm:"column:type;type:varchar(255);comment:类型"`
	Content string `gorm:"column:content;type:text;"`
	// 发生时间
	Ctime int64 `gorm:"column:ctime;comment:发生时间"`
}
