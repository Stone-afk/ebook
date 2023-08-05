package web

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	panic("")
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	panic("")
}

// SendSMSLoginCode 发送短信验证码
func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	panic("")
}

// LoginSMS 短信验证登录
func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	panic("")
}

// Login 用户登录接口
func (h *UserHandler) Login(ctx *gin.Context) {
	panic("")
}

// LoginJWT 用户登录接口
func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	panic("")
}

// SignUp 用户注册接口
func (h *UserHandler) SignUp(ctx *gin.Context) {
	panic("")
}

// Edit 用户编译信息
func (h *UserHandler) Edit(ctx *gin.Context) {
	panic("")
}

// Profile 用户详情
func (h *UserHandler) Profile(ctx *gin.Context) {
	panic("")
}

// ProfileJWT 用户详情, JWT 版本
func (h *UserHandler) ProfileJWT(ctx *gin.Context) {
	panic("")
}
