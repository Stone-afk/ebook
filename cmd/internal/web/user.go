package web

import (
	"ebook/cmd/internal/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	userIdKey = "userId"
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
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	// 当我们调用 Bind 方法的时候，如果有问题，Bind 方法已经直接写响应回去了
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或者密码不正确，请重试")
		return
	}
	sess := sessions.Default(ctx)
	sess.Set(userIdKey, u.Id)
	sess.Options(sessions.Options{
		// 60 秒过期
		MaxAge: 60,
	})
	err = sess.Save()
	if err != nil {
		ctx.String(http.StatusOK, "服务器异常")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
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
	type Profile struct {
		Email string
	}
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	u, err := h.svc.Profile(ctx, id)
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: u.Email,
	})
}

// ProfileJWT 用户详情, JWT 版本
func (h *UserHandler) ProfileJWT(ctx *gin.Context) {
	panic("")
}
