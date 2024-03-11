package handler

import articlev1 "ebook/cmd/api/proto/gen/article/v1"

type RewardReq struct {
	Id  int64 `json:"id"`
	Amt int64 `json:"amt"`
}

type LikeReq struct {
	Id int64 `json:"id"`
	// 点赞和取消点赞，我都准备复用这个
	Like bool `json:"like"`
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDTO(uid int64) *articlev1.Article {
	return &articlev1.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: &articlev1.Author{
			Id: uid,
		},
	}
}

type ArticleVO struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	Status  int32  `json:"status"`
	Author  string `json:"author"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
	// 计数
	ReadCnt    int64 `json:"read_cnt"`
	LikeCnt    int64 `json:"like_cnt"`
	CollectCnt int64 `json:"collect_cnt"`

	// 我个人有没有收藏，有没有点赞
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type ListReq struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

type CollectReq struct {
	Id  int64 `json:"id"`
	Cid int64 `json:"cid"`
}
