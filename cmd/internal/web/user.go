package web

import (
	"ebook/cmd/internal/service"
	"github.com/gin-gonic/gin"
)

var _ handler = &UserHandler{}

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 直接注册
	//server.POST("/users/signup", c.SignUp)
	//server.POST("/users/login", c.Login)
	//server.POST("/users/edit", c.Edit)
	//server.GET("/users/profile", c.Profile)

	// 分组注册
	ug := server.Group("/users")
	ug.POST("/signup", h.SignUp)
	// session 机制
	//ug.POST("/login", c.Login)
	// JWT 机制
	ug.POST("/login", h.LoginJWT)
	ug.POST("/edit", h.Edit)
	//ug.GET("/profile", c.Profile)
	ug.GET("/profile", h.ProfileJWT)
	ug.POST("/login_sms", h.LoginSMS)
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
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
