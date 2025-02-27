package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	myjwt "Webook/webook/internal/web/jwt"
	"Webook/webook/pkg/logger"
	"net/http"
	"time"

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

// ------------------------------------------------------------
// 查询部分
// ------------------------------------------------------------

// ArticlePage 文章分页
type ArticlePage struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ArticleVO struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     uint8  `json:"status"`
	Ctime      string `json:"ctime"`
	Utime      string `json:"utime"`
}

// List 获取文章列表
func (a *ArticleHandler) List(ctx *gin.Context) {
	var page ArticlePage
	if err := ctx.BindJSON(&page); err != nil {
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

	articles, err := a.svc.List(ctx, userId, page.Limit, page.Offset)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取文章列表失败",
			logger.Int64("limit", int64(page.Limit)),
			logger.Int64("offset", int64(page.Offset)),
			logger.Int64("userId", userId),
			logger.Error(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "获取文章列表成功",
		Data: toArticleVOs(articles),
	})
}

func toArticleVOs(arts []domain.Article) []ArticleVO {
	result := make([]ArticleVO, 0)
	for _, art := range arts {
		result = append(result, ArticleVO{
			Id:       art.Id,
			Title:    art.Title,
			Abstract: art.Abstract(),
			Status:   art.Status.ToUint8(),
			Ctime:    art.Ctime.Format(time.DateTime),
			Utime:    art.Utime.Format(time.DateTime),
		})
	}
	return result
}
