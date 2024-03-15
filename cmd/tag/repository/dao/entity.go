package dao

// Tag ID=1 => uid = 123
// ID = 1 => uid =234
type Tag struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 我要不要在这里创建一个唯一索引<uid, name>
	Name string `gorm:"type=varchar(4096)"`
	// 要在 uid 上创建一个索引
	// 因为你有一个典型的根据 uid 来查询的场景
	Uid   int64 `gorm:"index"`
	Ctime int64
	Utime int64
}

// TagBiz 某个人对某个资源打了标签
type TagBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	BizId int64  `gorm:"index:biz_type_id"`
	Biz   string `gorm:"index:biz_type_id"`
	// 冗余字段，加快查询和删除
	// 这个字段可以删除的
	Uid int64 `gorm:"index"`
	//TagName string
	Tid   int64
	Tag   *Tag  `gorm:"ForeignKey:Tid;AssociationForeignKey:Id;constraint:OnDelete:CASCADE"`
	Ctime int64 `bson:"ctime,omitempty"`
	Utime int64 `bson:"utime,omitempty"`
}
