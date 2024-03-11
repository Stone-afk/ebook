package handler

import (
	"ebook/cmd/pkg/ginx"
	"github.com/gin-gonic/gin"
)

type Result = ginx.Result[any]

type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type handler interface {
	RegisterRoutes(s *gin.Engine)
}
