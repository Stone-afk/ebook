package service

import (
	"context"
	"ebook/cmd/im/domain"
	"ebook/cmd/pkg/logger"
	"fmt"
	"github.com/ecodeclub/ekit/net/httpx"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type RESTUserService struct {
	// 部署 IM 时候配置的 IM Secret，默认是 openIM123
	secret   string
	baseHost string
	client   *http.Client
	log      logger.Logger
}

func (s *RESTUserService) Sync(ctx context.Context, user domain.User) error {
	spanCtx := trace.SpanContextFromContext(ctx)
	var traceId string
	if spanCtx.HasSpanID() {
		traceId = spanCtx.TraceID().String()
	} else {
		// 随便生成一个，但是这样链路就拼接不起来了
		traceId = uuid.New().String()
	}
	var resp response
	err := httpx.NewRequest(ctx,
		http.MethodPost, s.baseHost+"/user/user_register").JSONBody(syncUserRequest{
		Secret: s.secret,
		Users:  []domain.User{user},
	}).Client(s.client).AddHeader("operationID", traceId).
		Do().JSONScan(&resp)
	if err != nil {
		s.log.Error("同步用户到im系统异常", logger.Error(err))
		return err
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("同步数据失败 %d, %s, %s", resp.ErrCode, resp.ErrMsg, resp.ErrDlt)
	}
	return nil
}

func NewRESTUserService(secret Secret, baseHost BaseHost, client *http.Client, log logger.Logger) UserService {
	return &RESTUserService{
		log:      log,
		client:   client,
		secret:   string(secret),
		baseHost: string(baseHost),
	}
}

type syncUserRequest struct {
	Secret string        `json:"secret"`
	Users  []domain.User `json:"users"`
}

type response struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
}
