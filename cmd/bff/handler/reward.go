package handler

import (
	rewardv1 "ebook/cmd/api/proto/gen/reward/v1"
	"ebook/cmd/pkg/ginx"
	"ebook/cmd/pkg/logger"
	"github.com/gin-gonic/gin"
)

type RewardHandler struct {
	l            logger.Logger
	rewardClient rewardv1.RewardServiceClient
	//artClient articlev1.ArticleServiceClient
}

func NewRewardHandler(rewardClient rewardv1.RewardServiceClient) *RewardHandler {
	return &RewardHandler{rewardClient: rewardClient}
}

func (h *RewardHandler) RegisterRoutes(server *gin.Engine) {
	rg := server.Group("/reward")
	rg.POST("/detail", ginx.WrapClaimsAndReq[GetRewardReq](h.GetReward))
}

func (h *RewardHandler) GetReward(
	ctx *gin.Context,
	req GetRewardReq,
	claims ginx.UserClaims) (Result, error) {
	resp, err := h.rewardClient.GetReward(ctx, &rewardv1.GetRewardRequest{
		Rid: req.Rid,
		Uid: claims.Id,
	})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{
		// 暂时也就是只需要状态
		Data: resp.Status.String(),
	}, nil
}

type GetRewardReq struct {
	Rid int64
}

type RewardArticleReq struct {
	Aid int64 `json:"aid"`
	Amt int64 `json:"amt"`
}
