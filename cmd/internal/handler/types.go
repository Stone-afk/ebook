package handler

import "github.com/gin-gonic/gin"

type handler interface {
	RegisterRoutes(s *gin.Engine)
}
