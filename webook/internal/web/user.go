package web

import (
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
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

func NewUserHandler() *UserHandler {
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
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
	ctx.String(http.StatusOK, "user signup")
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {
	// 在页面上显示 hello world
	ctx.JSON(200, gin.H{ // 返回一个json格式的数据
		"message": "hello world",
	})

}