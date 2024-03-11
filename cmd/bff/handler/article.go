package handler

import (
	articlev1 "ebook/cmd/api/proto/gen/article/v1"
	intrv1 "ebook/cmd/api/proto/gen/intr/v1"
	rewardv1 "ebook/cmd/api/proto/gen/reward/v1"
	ijwt "ebook/cmd/bff/handler/jwt"
	"ebook/cmd/pkg/ginx"
	"ebook/cmd/pkg/logger"
	"errors"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

type ArticleHandler struct {
	svc            articlev1.ArticleServiceClient
	interactiveSvc intrv1.InteractiveServiceClient
	rewardSvc      rewardv1.RewardServiceClient
	l              logger.Logger
	biz            string
}

func NewArticleHandler(svc articlev1.ArticleServiceClient,
	interactiveSvc intrv1.InteractiveServiceClient,
	rewardSvc rewardv1.RewardServiceClient,
	l logger.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:            svc,
		rewardSvc:      rewardSvc,
		interactiveSvc: interactiveSvc,
		l:              l,
		biz:            "article",
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	// 修改
	//g.PUT("/")
	// 新增
	//g.POST("/")
	// g.DELETE("/a_id")
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	// Withdraw 撤回已发布贴子
	g.POST("/withdraw", h.Withdraw)
	// 在有 list 等路由的时候，无法这样注册
	// g.GET("/:id", a.Detail)
	g.GET("/detail/:id", ginx.WrapToken[ijwt.UserClaims](h.l, h.Detail))
	// 理论上来说应该用 GET的，但是我实在不耐烦处理类型转化
	// 直接 POST，JSON 转一了百了。
	g.POST("/list", ginx.WrapBodyAndToken[ListReq, ijwt.UserClaims](h.l, h.List))

	pub := g.Group("/pub")
	pub.GET("/:id", h.PubDetail)
	// 点赞是这个接口，取消点赞也是这个接口
	// RESTful 风格
	//pub.POST("/like/:id", ginx.WrapBodyAndToken[LikeReq,
	//	ijwt.UserClaims](h.Like))
	//pub.POST("/cancel_like", ginx.WrapBodyAndToken[LikeReq,
	//	ijwt.UserClaims](h.Like))
	pub.POST("/like", ginx.WrapBodyAndToken[LikeReq,
		ijwt.UserClaims](h.l, h.Like))
	pub.POST("/collect", ginx.WrapBodyAndToken[CollectReq,
		ijwt.UserClaims](h.l, h.Collect))
	// 打赏
	pub.POST("/reward", ginx.WrapClaimsAndReq[RewardReq](h.Reward))

}

func (h *ArticleHandler) Reward(
	ctx *gin.Context,
	req RewardReq,
	uc ijwt.UserClaims) (Result, error) {
	artResp, err := h.svc.GetPublishedById(ctx.Request.Context(),
		&articlev1.GetPublishedByIdRequest{
			Id: req.Id,
		})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	art := artResp.GetArticle()
	resp, err := h.rewardSvc.PreReward(ctx.Request.Context(), &rewardv1.PreRewardRequest{
		Biz:       "article",
		BizId:     art.Id,
		BizName:   art.Title,
		TargetUid: art.Author.GetId(),
		Uid:       uc.UserId,
		Amt:       req.Amt,
	})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{
		Data: map[string]any{
			"codeURL": resp.CodeUrl,
			"rid":     resp.Rid,
		},
	}, nil
}

func (h *ArticleHandler) Collect(ctx *gin.Context, req CollectReq, uc ijwt.UserClaims) (Result, error) {
	_, err := h.interactiveSvc.Collect(ctx, &intrv1.CollectRequest{
		Biz: h.biz, BizId: req.Id, Cid: req.Cid, Uid: uc.UserId,
	})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}

func (h *ArticleHandler) Like(ctx *gin.Context, req LikeReq, uc ijwt.UserClaims) (Result, error) {
	var err error
	if req.Like {
		_, err = h.interactiveSvc.Like(ctx, &intrv1.LikeRequest{
			Biz: h.biz, BizId: req.Id, Uid: uc.UserId,
		})
	} else {
		_, err = h.interactiveSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz: h.biz, BizId: req.Id, Uid: uc.UserId,
		})
	}
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		h.l.Error("输入的 ID 不对", logger.Error(err))
		return
	}
	uc := ctx.MustGet("users").(ijwt.UserClaims)
	// 使用 error group 来同时查询数据
	var (
		eg       errgroup.Group
		artResp  *articlev1.GetPublishedByIdResponse
		intrResp *intrv1.GetResponse
	)
	eg.Go(func() error {
		var er error
		artResp, er = h.svc.GetPublishedById(ctx, &articlev1.GetPublishedByIdRequest{
			Id: id, Uid: uc.UserId,
		})
		if er != nil {
			h.l.Error("获取发布文章详情失败", logger.Error(err))
		}
		return er
	})
	eg.Go(func() error {
		var er error
		// 这个地方可以容忍错误
		intrResp, er = h.interactiveSvc.Get(ctx, &intrv1.GetRequest{
			Biz: h.biz, BizId: id, Uid: uc.UserId,
		})
		if er != nil {
			h.l.Error("获取文章点赞收藏数量失败", logger.Error(err))
		}
		// 这种是容错的写法
		//if err != nil {
		//	// 记录日志
		//}
		//return nil
		return er
	})
	// 在这儿等，要保证前面两个
	err = eg.Wait()
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "系统错误",
		})
		return
	}

	// 增加阅读计数。
	//go func() {
	//	// 你都异步了，怎么还说有巨大的压力呢？
	//	// 开一个 goroutine，异步去执行
	//	er := h.intrSvc.IncrReadCnt(ctx, h.biz, art.Id)
	//	if er != nil {
	//		h.l.Error("增加阅读计数失败",
	//			logger.Int64("aid", art.Id),
	//			logger.Error(err))
	//	}
	//}()
	art := artResp.GetArticle()
	intr := intrResp.GetIntr()
	// 这个功能是不是可以让前端，主动发一个 HTTP 请求，来增加一个计数？
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status,
			Content: art.Content,
			// 要把作者信息带出去
			Author:     art.Author.Name,
			Ctime:      art.Ctime.AsTime().Format(time.DateTime),
			Utime:      art.Utime.AsTime().Format(time.DateTime),
			Liked:      intr.Liked,
			Collected:  intr.Collected,
			LikeCnt:    intr.LikeCnt,
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
		},
	})
}

func (h *ArticleHandler) Detail(ctx *gin.Context, uc ijwt.UserClaims) (Result, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		//ctx.JSON(http.StatusOK, )
		//a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return Result{
			Code: 4,
			Msg:  "参数错误",
		}, err
	}
	resp, err := h.svc.GetById(ctx, &articlev1.GetByIdRequest{Id: id})
	if err != nil {
		//ctx.JSON(http.StatusOK, )
		//a.l.Error("获得文章信息失败", logger.Error(err))
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	art := resp.GetArticle()
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != uc.UserId {
		//ctx.JSON(http.StatusOK)
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		//a.l.Error("非法访问文章，创作者 ID 不匹配",
		//	logger.Int64("uid", usr.Id))
		return Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		}, fmt.Errorf("非法访问文章，创作者 ID 不匹配 %d", uc.UserId)
	}
	return Result{
		Data: ArticleVO{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status,
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: art.Ctime.AsTime().Format(time.DateTime),
			Utime: art.Utime.AsTime().Format(time.DateTime),
		},
	}, nil
}

func (h *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (Result, error) {
	// 对于批量接口来说，要小心批次大小
	if req.Limit > 100 {
		return Result{
			Code: 4,
			// 我会倾向于不告诉前端批次太大
			// 因为一般你和前端一起完成任务的时候
			// 你们是协商好了的，所以会进来这个分支
			// 就表明是有人跟你过不去
			Msg: "请求有误",
		}, errors.New("单次请求超出最大数量限制")
	}
	resp, err := h.svc.List(ctx, &articlev1.ListRequest{Author: uc.UserId,
		Offset: req.Offset, Limit: req.Limit})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{
		Data: slice.Map[*articlev1.Article, ArticleVO](resp.Articles,
			func(idx int, src *articlev1.Article) ArticleVO {
				return ArticleVO{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract,
					Status:   src.Status,
					// 这个列表请求，不需要返回内容
					//Content: src.Content,
					// 这个是创作者看自己的文章列表，也不需要这个字段
					//Author: src.Author
					Ctime: src.Ctime.AsTime().Format(time.DateTime),
					Utime: src.Utime.AsTime().Format(time.DateTime),
				}
			}),
	}, nil
}

//func (h *ArticleHandler) List(ctx *gin.Context) {
//	type Req struct {
//		Offset int `json:"offset"`
//		Limit  int `json:"limit"`
//	}
//	var req Req
//	if err := ctx.Bind(&req); err != nil {
//		h.l.Error("反序列化请求失败", logger.Error(err))
//		return
//	}
//	// 对于批量接口来说，要小心批次大小
//	if req.Limit > 100 {
//		ctx.JSON(http.StatusOK, Result{
//			Code: 4,
//			// 我会倾向于不告诉前端批次太大
//			// 因为一般你和前端一起完成任务的时候
//			// 你们是协商好了的，所以会进来这个分支
//			// 就表明是有人跟你过不去
//			Msg: "请求有误",
//		})
//		h.l.Error("获得用户会话信息失败")
//		return
//	}
//	usr, ok := ctx.MustGet("user").(ijwt.UserClaims)
//	if !ok {
//		ctx.JSON(http.StatusOK, Result{
//			Code: 5,
//			Msg:  "系统错误",
//		})
//		h.l.Error("获得用户会话信息失败")
//		return
//	}
//	arts, err := h.svc.List(ctx, usr.UserId, req.Offset, req.Limit)
//	if err != nil {
//		ctx.JSON(http.StatusOK, Result{
//			Code: 5,
//			Msg:  "系统错误",
//		})
//		h.l.Error("获得用户会话信息失败")
//		return
//	}
//	ctx.JSON(http.StatusOK, Result{
//		Data: slice.Map[domain.Article, ArticleVo](arts,
//			func(idx int, src domain.Article) ArticleVo {
//				return ArticleVo{
//					Id:       src.Id,
//					Title:    src.Title,
//					Abstract: src.Abstract(),
//					Status:   src.Status.ToUint8(),
//					// 这个列表请求，不需要返回内容
//					//Content: src.Content,
//					// 这个是创作者看自己的文章列表，也不需要这个字段
//					//Author: src.Author
//					Ctime: src.Ctime.Format(time.DateTime),
//					Utime: src.Utime.Format(time.DateTime),
//				}
//			}),
//	})
//}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
		return
	}
	// 检测输入，跳过这一步
	// 调用 svc 的代码
	_, err := h.svc.Withdraw(ctx, &articlev1.WithdrawRequest{
		Uid: claims.UserId, Id: req.Id})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		h.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	userClaims, ok := ctx.MustGet("user").(ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获得用户会话信息失败")
		return
	}
	resp, err := h.svc.Publish(ctx, &articlev1.PublishRequest{
		Article: &articlev1.Article{
			Id:      req.Id,
			Title:   req.Title,
			Content: req.Content,
			Author: &articlev1.Author{
				Id: userClaims.UserId,
			},
		},
	})
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("发布文章失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: resp.Id,
	})
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
		return
	}
	// 检测输入，跳过这一步
	// 调用 svc 的代码
	id, err := h.svc.Save(ctx, &articlev1.SaveRequest{Article: req.toDTO(claims.UserId)})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}
