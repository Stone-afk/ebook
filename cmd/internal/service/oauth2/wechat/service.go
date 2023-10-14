package wechat

import (
	"context"
	"ebook/cmd/internal/domain"
	"ebook/cmd/internal/service/oauth2"
	"ebook/cmd/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

var _ oauth2.Service = &Service{}

type Service struct {
	appId     string
	appSecret string
	client    *http.Client
	//cmd       redis.Cmdable
	l logger.Logger
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenID  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}

func NewService(appId string, appSecret string, client *http.Client, l logger.Logger) oauth2.Service {
	return &Service{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
		l:         l,
	}
}

func (s *Service) AuthURL(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redirect"
	// 如果在这里存 state，假如说我存 redis
	//s.cmd.Set(ctx, "my-state"+state, state, time.Minute)
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

func (s *Service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	//const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	//target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	// 这是另外一种写法
	const baseURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	queryParams := url.Values{}
	queryParams.Set("appid", s.appId)
	queryParams.Set("secret", s.appSecret)
	queryParams.Set("code", code)
	queryParams.Set("grant_type", "authorization_code")
	target := baseURL + "?" + queryParams.Encode()

	//resp, err := http.Get(target)
	//req, err := http.NewRequest(http.MethodGet, target, nil)
	// 会产生复制，性能极差，比如说你的 URL 很长
	//req = req.WithContext(ctx)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	// 整个响应都读出来，不推荐，因为 Unmarshal 再读一遍，合计两遍
	//body, err := io.ReadAll(resp.Body)
	//err = json.Unmarshal(body, &res)

	// 只读一遍
	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误响应，错误码：%d，错误信息：%s", res.ErrCode, res.ErrMsg)
	}

	// 攻击者的 state
	//str := s.cmd.Get(ctx, "my-state"+state).String()
	//if str != state {
	//	// 不相等
	//}

	return domain.WechatInfo{
		OpenID:  res.OpenID,
		UnionID: res.UnionID,
	}, nil
}
