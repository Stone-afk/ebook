package domain

type SearchResult struct {
	Users    []User
	Articles []Article
	BizTags  []BizTags
}

type SearchUserResult struct {
	Users []User
}

type SearchArticleResult struct {
	Articles []Article
}

type SearchBizTagsResult struct {
	BizTags []BizTags
}
