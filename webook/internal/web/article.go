package web

import (
	"Webook/webook/internal/domain"
	"Webook/webook/internal/service"
	myjwt "Webook/webook/internal/web/jwt"
	"Webook/webook/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	svc    service.ArticleService
	logger logger.Logger

	// 阅读，点赞，收藏
	biz      string
	interSvc service.InteractiveService
}

type ArticleReaderHandler struct {
	svc    service.ArticleService
	logger logger.Logger

	// 阅读，点赞，收藏
	biz      string
	interSvc service.InteractiveService
	rankSvc  service.RankingService
	userSvc  service.UserService
}

func NewArticleHandler(svc service.ArticleService, interSvc service.InteractiveService, logger logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:      svc,
		logger:   logger,
		biz:      "article",
		interSvc: interSvc,
	}
}

func NewArticleReaderHandler(svc service.ArticleService, interSvc service.InteractiveService, rankSvc service.RankingService, userSvc service.UserService, logger logger.Logger) *ArticleReaderHandler {
	return &ArticleReaderHandler{
		svc:      svc,
		logger:   logger,
		biz:      "article",
		interSvc: interSvc,
		rankSvc:  rankSvc,
		userSvc:  userSvc,
	}
}

func (a *ArticleHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.POST("/edit", a.Edit)
	ug.POST("/publish", a.Publish)
	ug.POST("/withdraw", a.Withdraw)
	ug.POST("/delete", a.Delete)

	// 文章列表
	ug.POST("/list", a.List)
	// 文章详情
	ug.GET("/detail/:id", a.Detail)

}

func (a *ArticleReaderHandler) RegisterRoutes(ug *gin.RouterGroup) {
	ug.GET("/:id", a.PublicDetail)
	ug.POST("/like", a.Like)
	ug.POST("/collect", a.Collect)
	ug.POST("/rank/list", a.RankingList)
	ug.POST("/list", a.PublicList)
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

	// 阅读量，点赞数，收藏数
	ReadCnt    int64 `json:"readCnt"`
	LikeCnt    int64 `json:"likeCnt"`
	CollectCnt int64 `json:"collectCnt"`
	Liked      bool  `json:"liked"`
	Collected  bool  `json:"collected"`
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

	bizIds := slice.Map(articles, func(idx int, art domain.Article) int64 {
		return art.Id
	})

	interMap, err := a.interSvc.GetInterMapByBizIds(ctx, a.biz, bizIds, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取文章交互信息失败",
			logger.Error(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "获取文章列表成功",
		Data: toArticleVOs(articles, interMap, nil),
	})
}

func toArticleVOs(arts []domain.Article, interMap map[int64]domain.Interactive, authorMap map[int64]string) []ArticleVO {
	result := make([]ArticleVO, 0)
	for _, art := range arts {
		result = append(result, ArticleVO{
			Id:         art.Id,
			Title:      art.Title,
			Abstract:   art.Abstract(),
			Status:     art.Status.ToUint8(),
			Ctime:      art.Ctime.Format(time.DateTime),
			Utime:      art.Utime.Format(time.DateTime),
			AuthorId:   art.Author.Id,
			AuthorName: authorMap[art.Author.Id],
			ReadCnt:    interMap[art.Id].ReadCnt,
			LikeCnt:    interMap[art.Id].LikeCnt,
			CollectCnt: interMap[art.Id].CollectCnt,
			Liked:      interMap[art.Id].Liked,
			Collected:  interMap[art.Id].Collected,
		})
	}
	return result
}

// Detail 获取文章详情
func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "id 格式错误",
		})
		a.logger.Error("id 格式错误",
			logger.String("id", idStr),
			logger.Error(err),
		)
		return
	}

	article, err := a.svc.Detail(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取文章详情失败",
			logger.Int64("id", id),
			logger.Error(err),
		)
		return
	}

	// 检查用户是否是作者
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

	if article.Author.Id != userId {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "您无权限编辑其他用户的文章",
		})
		a.logger.Error("用户无权限编辑其他人的文章",
			logger.Int64("articleId", article.Id),
			logger.Int64("userId", userId),
		)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "获取文章详情成功",
		Data: ArticleVO{
			Id:         article.Id,
			Title:      article.Title,
			Content:    article.Content,
			AuthorId:   article.Author.Id,
			AuthorName: article.Author.Name,
			Status:     article.Status.ToUint8(),
			Ctime:      article.Ctime.Format(time.DateTime),
			Utime:      article.Utime.Format(time.DateTime),
		},
	})
}

func (a *ArticleReaderHandler) PublicDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "id 格式错误",
		})
		a.logger.Error("id 格式错误",
			logger.String("id", idStr),
			logger.Error(err),
		)
		return
	}

	article, err := a.svc.PublicDetail(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取Public文章详情失败",
			logger.Int64("id", id),
			logger.Error(err),
		)
		return
	}

	// 增加阅读计数，保证 Interactive 中存在 (biz,bizId) 的记录
	// IncreaseReadCnt 采用 Upsert 的写法，如果记录不存在，则创建，如果记录存在，则更新
	err = a.interSvc.IncreaseReadCnt(ctx, a.biz, article.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("增加阅读计数失败",
			logger.Int64("id", article.Id),
			logger.Error(err),
		)
		return
	}

	userClaims := ctx.MustGet("claims").(*myjwt.UserClaims)
	interactive, err := a.interSvc.Get(ctx, a.biz, article.Id, userClaims.UserId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取互动信息失败",
			logger.Int64("id", article.Id),
			logger.Error(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "获取Public文章详情成功",
		Data: ArticleVO{
			Id:         article.Id,
			Title:      article.Title,
			Content:    article.Content,
			AuthorId:   article.Author.Id,
			AuthorName: article.Author.Name,
			Status:     article.Status.ToUint8(),
			Ctime:      article.Ctime.Format(time.DateTime),
			Utime:      article.Utime.Format(time.DateTime),

			ReadCnt:    interactive.ReadCnt,
			LikeCnt:    interactive.LikeCnt,
			CollectCnt: interactive.CollectCnt,
			Liked:      interactive.Liked,
			Collected:  interactive.Collected,
		},
	})
}

func (a *ArticleReaderHandler) Like(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
		// 点赞还是取消点赞
		Like bool `json:"like"`
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

	var err error
	if req.Like {
		err = a.interSvc.IncreaseLike(ctx, a.biz, req.Id, userId)
	} else {
		err = a.interSvc.DecreaseLike(ctx, a.biz, req.Id, userId)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("点赞/取消点赞失败",
			logger.Int64("articleId", req.Id),
			logger.Int64("userId", userId),
			logger.Error(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "点赞/取消点赞成功",
	})
}

func (a *ArticleReaderHandler) Collect(ctx *gin.Context) {
	type Req struct {
		Id  int64 `json:"id"`
		Cid int64 `json:"cid"`
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

	err := a.interSvc.Collect(ctx, a.biz, req.Id, req.Cid, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("收藏失败",
			logger.Int64("articleId", req.Id),
			logger.Int64("collectionId", req.Cid),
			logger.Int64("userId", userId),
			logger.Error(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "收藏成功",
	})
}

func (a *ArticleReaderHandler) RankingList(ctx *gin.Context) {
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

	articles, err := a.rankSvc.GetFromCache(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取榜单列表失败",
			logger.Int64("limit", int64(page.Limit)),
			logger.Int64("offset", int64(page.Offset)),
			logger.Int64("userId", userId),
			logger.Error(err),
		)
		return
	}

	bizIds := slice.Map(articles, func(idx int, art domain.Article) int64 {
		return art.Id
	})

	// 获取文章点赞数
	interMap, err := a.interSvc.GetInterMapByBizIds(ctx, "article", bizIds, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取文章交互信息失败",
			logger.Error(err),
		)
		return
	}

	// 获取作者信息
	userIds := slice.Map(articles, func(idx int, art domain.Article) int64 {
		return art.Author.Id
	})
	authorMap, err := a.userSvc.GetNameMapByIds(ctx, userIds)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取作者信息失败",
			logger.Error(err),
		)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "获取文章列表成功",
		Data: toArticleVOs(articles, interMap, authorMap),
	})
}

func (a *ArticleReaderHandler) PublicList(ctx *gin.Context) {
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
	now := time.Now()
	articles, err := a.svc.PublicList(ctx, now, page.Offset, page.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取榜单列表失败",
			logger.Int64("limit", int64(page.Limit)),
			logger.Int64("offset", int64(page.Offset)),
			logger.Int64("userId", userId),
			logger.Error(err),
		)
		return
	}

	// 获取文章点赞数
	bizIds := slice.Map(articles, func(idx int, art domain.Article) int64 {
		return art.Id
	})
	interMap, err := a.interSvc.GetInterMapByBizIds(ctx, "article", bizIds, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取文章交互信息失败",
			logger.Error(err),
		)
		return
	}

	// 获取作者信息
	userIds := slice.Map(articles, func(idx int, art domain.Article) int64 {
		return art.Author.Id
	})
	authorMap, err := a.userSvc.GetNameMapByIds(ctx, userIds)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("获取作者信息失败",
			logger.Error(err),
		)
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "获取文章列表成功",
		Data: toArticleVOs(articles, interMap, authorMap),
	})
}
