package handler

import (
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	userIdKey            = "userId"
)

var _ handler = &UserHandler{}

type UserHandler struct {
	svc              service.UserService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
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
	h.setJWTToken(ctx, u.Id)
	ctx.String(http.StatusOK, "登录成功")
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, UserClaims{
		Id:        uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			// 演示目的设置为一分钟过期
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			// 在压测的时候，要将过期时间设置更长一些
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	})
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
}

// SignUp 用户注册接口
func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	// 当我们调用 Bind 方法的时候，如果有问题，Bind 方法已经直接写响应回去了
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := h.emailRegexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "邮箱不正确")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	isPassword, err := h.passwordRegexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含数字、特殊字符，并且长度不能小于 8 位")
		return
	}

	err = h.svc.Signup(ctx.Request.Context(),
		domain.User{Email: req.Email, Password: req.ConfirmPassword})

	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "该邮箱已被注册")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "注册失败，服务器异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
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
	type Profile struct {
		Email string
		Phone string
	}
	uc := ctx.MustGet("user").(UserClaims)
	u, err := h.svc.Profile(ctx, uc.Id)
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