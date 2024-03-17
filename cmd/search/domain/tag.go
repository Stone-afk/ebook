package domain

type Tag struct {
	Uid   int64
	Biz   string
	BizId int64
	// 只传递 string
	Tags []string
}
