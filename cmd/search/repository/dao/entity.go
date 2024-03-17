package dao

const ArticleIndexName = "article_index"

type Article struct {
	Id      int64    `json:"id"`
	Title   string   `json:"title"`
	Status  int32    `json:"status"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

const UserIndexName = "user_index"

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
}

const TagIndexName = "tags_index"

type BizTags struct {
	Uid   int64  `json:"uid"`
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`
	// 只传递 string
	Tags []string `json:"tags"`
}
