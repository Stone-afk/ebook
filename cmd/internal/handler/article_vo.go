package handler

import "ebook/cmd/internal/domain"

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleVo struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	// 摘要
	Abstract string `json:"abstract"`
	// 内容
	Content string `json:"content"`
	Status  uint8  `json:"status"`
	Author  string `json:"author"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}

func (req ArticleReq) toDomain(userId int64) domain.Article {
	return domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: userId,
		},
	}
}
