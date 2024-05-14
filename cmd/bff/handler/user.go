package handler

import (
	codev1 "ebook/cmd/api/proto/gen/code/v1"
	userv1 "ebook/cmd/api/proto/gen/user/v1"
	ijwt "ebook/cmd/bff/handler/jwt"
	"ebook/cmd/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`

	userIdKey = "userId"
	bizLogin  = "login"
)

type SignUpReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

var _ handler = &UserHandler{}

type UserHandler struct {
	svc              userv1.UserServiceClient
	codeSvc          codev1.CodeServiceClient
	emailRegexExp    *regexp.Regexp
	passwordRegexExp *regexp.Regexp
	ijwt.Handler
}

func (h *UserHandler) RegisterRoutes(s *gin.Engine) {
	//TODO implement me
	panic("implement me")
}

func (h *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := h.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录OK",
	})
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	// 假定长 token 也放在这里
	tokenStr := h.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(tokenStr, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RefreshTokenKey, nil
	})
	// 这边要保持和登录校验一直的逻辑，即返回 401 响应
	if err != nil || token == nil || !token.Valid {
		ctx.JSON(http.StatusUnauthorized, Result{Code: 4, Msg: "请登录"})
		return
	}
	// 校验 ssid
	err = h.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 系统错误或者用户已经主动退出登录了
		// 这里也可以考虑说，如果在 Redis 已经崩溃的时候，
		// 就不要去校验是不是已经主动退出登录了。
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = h.SetJWTToken(ctx, rc.UserId, rc.Ssid)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, Result{Code: 4, Msg: "请登录"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "刷新成功"})
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
	_, err := h.codeSvc.Send(ctx, &codev1.CodeSendRequest{
		Biz: bizLogin, Phone: req.Phone,
	})
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
	resp, err := h.codeSvc.Verify(ctx, &codev1.VerifyRequest{
		Biz: bizLogin, Phone: req.Phone, InputCode: req.Code,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("用户手机号码登录失败", zap.Error(err))
		return
	}
	if resp.Answer {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	// 我这个手机号，会不会是一个新用户呢？
	// 这样子
	user, err := h.svc.FindOrCreate(ctx, &userv1.FindOrCreateRequest{
		Phone: req.Phone,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 用 uuid 来标识这一次会话
	// 短信登录使用长 token
	ssid := uuid.New().String()
	if err = h.SetJWTToken(ctx, user.User.Id, ssid); err != nil {
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
	u, err := h.svc.Login(ctx.Request.Context(), &userv1.LoginRequest{
		Email: req.Email, Password: req.Password})
	// TODO 利用 grpc 来传递错误码
	//if err == service.ErrInvalidUserOrPassword {
	//	ctx.String(http.StatusOK, "用户名或者密码不正确，请重试")
	//	return
	//}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	sess := sessions.Default(ctx)
	sess.Set(userIdKey, u.User.Id)
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
	u, err := h.svc.Login(ctx.Request.Context(), &userv1.LoginRequest{
		Email: req.Email, Password: req.Password,
	})
	// TODO 利用 grpc 来传递错误码
	//if err == service.ErrInvalidUserOrPassword {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: errs.UserInvalidOrPassword,
	//		Msg:  "用户名或者密码不正确，请重试",
	//	})
	//	return
	//}

	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if err = h.SetLoginToken(ctx, u.User.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
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

	_, err = h.svc.Signup(
		ctx.Request.Context(),
		&userv1.SignupRequest{User: &userv1.User{Email: req.Email, Password: req.ConfirmPassword}})

	// TODO 利用 grpc 来传递错误码
	//if err == service.ErrUserDuplicateEmail {
	//	ctx.String(http.StatusOK, "该邮箱已被注册")
	//	return
	//}

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
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	_, err = h.svc.UpdateNonSensitiveInfo(ctx,
		&userv1.UpdateNonSensitiveInfoRequest{
			User: &userv1.User{
				Id:       uc.UserId,
				Nickname: req.Nickname,
				AboutMe:  req.AboutMe,
				Birthday: timestamppb.New(birthday),
			},
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
	u, err := h.svc.Profile(ctx, &userv1.ProfileRequest{
		Id: id,
	})
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: u.User.Email,
	})
}

// ProfileJWT 用户详情, JWT 版本
func (h *UserHandler) ProfileJWT(ctx *gin.Context) {
	type Profile struct {
		Email string
		Phone string
	}
	uc := ctx.MustGet("user").(ijwt.UserClaims)
	resp, err := h.svc.Profile(ctx, &userv1.ProfileRequest{Id: uc.UserId})
	if err != nil {
		// 按照道理来说，这边 id 对应的数据肯定存在，所以要是没找到，
		// 那就说明是系统出了问题。
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.JSON(http.StatusOK, Profile{
		Email: resp.User.Email,
	})
}

func NewUserHandler(svc userv1.UserServiceClient,
	codeSvc codev1.CodeServiceClient, jwthdl ijwt.Handler) *UserHandler {
	return &UserHandler{
		svc:              svc,
		codeSvc:          codeSvc,
		emailRegexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		Handler:          jwthdl,
	}
}
