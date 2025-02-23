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
	ug.POST("/withdraw", a.Withdraw)
	ug.POST("/delete", a.Delete)
	ug.POST("/list", a.List)
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

// Withdraw 撤回文章
func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
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

	id, err := a.svc.Withdraw(ctx, domain.Article{
		Id:     req.Id,
		Author: domain.Author{Id: userId},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "撤回成功",

		Data: id,
	})
}

// Delete 删除文章
func (a *ArticleHandler) Delete(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
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

	id, err := a.svc.Delete(ctx, domain.Article{
		Id:     req.Id,
		Author: domain.Author{Id: userId},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "删除成功",
		Data: id,
	})
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "获取成功",
		Data: []domain.Article{
			{
				Id:      1,
				Title:   "标题1",
				Content: "内容1",
				Author:  domain.Author{Id: 1},
			},
			{
				Id:      2,
				Title:   "标题2",
				Content: "内容2",
				Author:  domain.Author{Id: 2},
			},
		},
	})
}
