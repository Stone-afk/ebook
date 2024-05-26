package handler

import (
	"ebook/cmd/pkg/ginx"
	"github.com/gin-gonic/gin"
)

var (
	SessAuthKey       = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	SessEncryptionKey = []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3")
)

type Result = ginx.Result[any]

type Page struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type handler interface {
	RegisterRoutes(s *gin.Engine)
}
