package handler

import "github.com/gin-gonic/gin"

var (
	SessAuthKey       = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	SessEncryptionKey = []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3")
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type handler interface {
	RegisterRoutes(s *gin.Engine)
}
