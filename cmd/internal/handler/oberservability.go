package handler

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type ObservabilityHandler struct{}

func NewObservabilityHandler() *ObservabilityHandler {
	return &ObservabilityHandler{}
}

func (h *ObservabilityHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("test")
	g.GET("/random", h.Random)
}

func (h *ObservabilityHandler) Random(ctx *gin.Context) {
	num := rand.Int31n(1000)
	// 模拟响应时间
	time.Sleep(time.Millisecond * time.Duration(num))
	ctx.String(http.StatusOK, "OK")
}
