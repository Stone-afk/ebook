package wechat

import (
	"context"
	"ebook/cmd/oauth2/domain/wechat"
	"ebook/cmd/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	appId     string
	appSecret string
	client    *http.Client
	logger    logger.Logger
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(authURLPattern, s.appId, redirectURL, state), nil
}

func (s *service) VerifyCode(ctx context.Context, code string) (wechat.WechatInfo, error) {
	const baseURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	// 这是另外一种写法
	queryParams := url.Values{}
	queryParams.Set("appid", s.appId)
	queryParams.Set("secret", s.appSecret)
	queryParams.Set("code", code)
	queryParams.Set("grant_type", "authorization_code")
	accessTokenURL := baseURL + "?" + queryParams.Encode()

	req, err := http.NewRequest("GET", accessTokenURL, nil)
	if err != nil {
		return wechat.WechatInfo{}, err
	}

	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return wechat.WechatInfo{}, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	var res Result

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return wechat.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return wechat.WechatInfo{}, errors.New("换取 access_token 失败")
	}
	return wechat.WechatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionId,
	}, nil
}

func NewService(appId, appSecret string,
	logger logger.Logger) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		logger:    logger,
		client:    http.DefaultClient,
	}
}
