package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	"errors"
	"net/http"
	"time"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

// 注册路由
func (u *UserHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.POST("/signup", u.SignUp)
	// ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.EditJWT)
	// ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/login_sms/code/send", u.LoginSMSCodeSend)
	ug.POST("/login_sms", u.LoginSMSCodeVerify)
}

const (
	emailRegexPattern    = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	passwordRegexPattern = "^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>\\/?]).{8,}$"
)

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		codeSvc:     codeSvc,
	}
}

type SignUpReq struct {
	Email           string `json:"email"`
	ConfirmPassword string `json:"confirmPassword"`
	Password        string `json:"password"`
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	var req SignUpReq

	// Bind 会按照 Content-Type 来解析上下文中的数据到 req 中
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 邮箱格式校验
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}

	// 密码格式校验
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须包含至少一个字母、数字、特殊字符，长度至少8位")
		return
	}

	// 确认密码和密码一致
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	// 调用 service 层进行注册
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicate {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "user signup")
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserHandler) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 调用 service 层进行登录
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 登录成功，设置 session
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		// 设置 session 的过期时间
		// MaxAge: 60 * 30, // 30 min
		MaxAge: 60,
	})
	sess.Set("userId", user.Id)

	sess.Save()

	// 获取 session
	userId := sess.Get("userId")
	ctx.String(http.StatusOK, "登录成功，userId: %d", userId)
}

type UserClaims struct {
	UserId int64
	jwt.RegisteredClaims
	UserAgent string
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) {
	// claims 中存储用户的信息
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置 token 的过期时间: 1 分钟（和 lua 代码中的过期时间一致）
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserId:    uid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("sUwYXfLAdddhd1hyWJkWMd4gqQiFznp6"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
	//ctx.String(http.StatusOK, "JWT 登录成功，token: %s", tokenStr)
}
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 调用 service 层进行登录
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 设置 JWT token，保持登录状态
		u.setJWTToken(ctx, user.Id)
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}

}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录")
}

type EditReq struct {
	Nickname string `json:"nickname"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe"`
}

func (u *UserHandler) Edit(ctx *gin.Context) {
}
func (u *UserHandler) EditJWT(ctx *gin.Context) {
	print("edit jwt")
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 获取 JWT 中的用户信息
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	userClaims := claims.(*UserClaims)
	userId := userClaims.UserId

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "生日格式错误")
		return
	}

	// 调用 service 层进行编辑
	user := domain.User{
		Id:       userId,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	}
	err = u.svc.Edit(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{Msg: "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "编辑成功"})
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	// 在页面上显示 hello world
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")
	user, err := u.svc.Profile(ctx, userId.(int64))
	if err != nil {
		ctx.String(http.StatusBadRequest, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "hello world, %+v", user)

}

type ProfileJWTResp struct {
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Nickname string `json:"Nickname"`
	Birthday string `json:"Birthday"`
	AboutMe  string `json:"AboutMe"`
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	// 获取 JWT 中的用户信息
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	userClaims := claims.(*UserClaims)

	// 获取用户信息
	userId := userClaims.UserId
	user, err := u.svc.Profile(ctx, userId)
	if err != nil {
		ctx.String(http.StatusBadRequest, "系统错误")
		return
	}

	// 返回用户信息
	ctx.JSON(http.StatusOK,
		&ProfileJWTResp{
			Email:    user.Email,
			Phone:    user.Phone,
			Nickname: user.Nickname,
			Birthday: user.Birthday.Format(time.DateOnly),
			AboutMe:  user.AboutMe,
		},
	)
}

type LoginSMSCodeSendReq struct {
	Phone string `json:"phone"`
}

func (u *UserHandler) LoginSMSCodeSend(ctx *gin.Context) {
	var req LoginSMSCodeSendReq

	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 校验手机号是否合法
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机号不能为空",
		})
		return
	}

	err := u.codeSvc.Send(ctx, "login", req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{Msg: "发送成功"})
	case service.ErrCodeSendTooFrequent:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "发送验证码过于频繁，请稍后再试",
		})
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

type LoginSMSCodeVerifyReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

func (u *UserHandler) LoginSMSCodeVerify(ctx *gin.Context) {
	var req LoginSMSCodeVerifyReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.codeSvc.Verify(ctx, "login", req.Phone, req.Code)
	if err != nil {
		if errors.Is(err, service.ErrCodeVerifyTooManyTimes) {
			ctx.JSON(http.StatusOK, Result{
				Code: 5,
				Msg:  "验证码错误次数过多，请稍后再试",
			})
			return
		}
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误，请重新输入",
		})
		return
	}

	// 查找或创建用户
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 配置 JWT token，保持登录状态
	u.setJWTToken(ctx, user.Id)
	ctx.JSON(http.StatusOK, Result{
		Msg:  "验证码校验通过",
		Data: user,
	})
}
