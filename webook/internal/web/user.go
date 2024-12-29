package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

// 注册路由
func (u *UserHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.PUT("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

const (
	emailRegexPattern    = "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$"
	passwordRegexPattern = "^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[!@#$%^&*()_+\\-=\\[\\]{};':\"\\\\|,.<>\\/?]).{8,}$"
)

func NewUserHandler(svc *service.UserService) *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
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
	if err == service.ErrUserDuplicateEmail {
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

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录")
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	// 在页面上显示 hello world
	sess := sessions.Default(ctx)
	userId := sess.Get("userId")

	ctx.String(http.StatusOK, "hello world, userId: %d", userId)

}
