package domain

const (
	// ArticleStatusUnknown 为了避免零值之类的问题
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

// Article 可以同时表达线上库和制作库的概念吗？
// 可以同时表达，作者眼中的 Article 和读者眼中的 Article 吗？
type Article struct {
	Id      int64
	Title   string
	Content string
	// Author 要从用户来
	Author Author
	Status ArticleStatus
}

// Author 在帖子这个领域内，是一个值对象
type Author struct {
	Id   int64
	Name string
}

type ArticleStatus uint8

// ArticleStatusV1 如果你的状态很复杂，有很多行为（就是你要搞很多方法），状态里面需要一些额外字段
// 就用这个版本
type ArticleStatusV1 struct {
	Val  uint8
	Name string
}

var (
	ArticleStatusV1Unknown = ArticleStatusV1{Val: 0, Name: "unknown"}
)

type ArticleStatusV2 string