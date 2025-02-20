package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	myjwt "Webook/webook/internal/web/jwt"
	"Webook/webook/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.Logger
}

func NewArticleHandler(svc service.ArticleService, logger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    svc,
		logger: logger,
	}
}

func (a *ArticleHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.POST("/edit", a.Edit)
	ug.POST("/publish", a.Publish)
}

// Edit 编辑文章
type ArticleRequest struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	// 获取 JWT 中的用户信息
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取 JWT 中的用户信息失败")
		return
	}
	userClaims := claims.(*myjwt.UserClaims)
	userId := userClaims.UserId

	id, err := a.svc.Save(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: userId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("编辑文章失败", logger.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg:  "编辑成功",
		Data: id,
	})
}

// Publish 发布文章
func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleRequest
	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	// 获取 JWT 中的用户信息
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取 JWT 中的用户信息失败")
		return
	}
	userClaims := claims.(*myjwt.UserClaims)
	userId := userClaims.UserId

	id, err := a.svc.Publish(ctx, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{Id: userId},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("发布文章失败", logger.Error(err))
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "发布成功",
		Data: id,
	})
}
