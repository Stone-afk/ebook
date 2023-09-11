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
	biz               = "login"
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	userIdKey            = "userId"
)

var _ handler = &UserHandler{}

type UserHandler struct {
	svc              service.UserService
	codeSvc          service.CodeService
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	return &UserHandler{
		svc:              svc,
		codeSvc:          codeSvc,
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
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 是不是一个合法的手机号码
	// 考虑正则表达式
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "输入有误",
		})
		return
	}
	err := h.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

// LoginSMS 短信验证登录
func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这边，可以加上各种校验
	ok, err := h.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err == service.ErrCodeVerifyTooManyTimes {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证操作过于频繁",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	// 我这个手机号，会不会是一个新用户呢？
	// 这样子
	user, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 这边要怎么办呢？ 从哪来？
	if err = h.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过, 登录成功",
	})
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
	if err = h.setJWTToken(ctx, u.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
}

func (h *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
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
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
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
		ctx.String(http.StatusOK, "邮箱格式不正确")
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
		ctx.String(http.StatusOK, "注册失败，系统异常")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
}

// Edit 用户编译信息
func (h *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		// 注意，其它字段，尤其是密码、邮箱和手机，
		// 修改都要通过别的手段
		// 邮箱和手机都要验证
		// 密码更加不用多说了
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你可以尝试在这里校验。
	// 比如说你可以要求 Nickname 必须不为空
	// 校验规则取决于产品经理
	if req.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "昵称不能为空"})
	}

	if len(req.AboutMe) > 1024 {
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "关于我的信息过长"})
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		// 也就是说，我们其实并没有直接校验具体的格式
		// 而是如果你能转化过来，那就说明没问题
		ctx.JSON(http.StatusOK, Result{Code: 4, Msg: "日期格式不对"})
		return
	}
	uc := ctx.MustGet("user").(UserClaims)
	err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:       uc.Id,
		Nickname: req.Nickname,
		AboutMe:  req.AboutMe,
		Birthday: birthday,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})

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
