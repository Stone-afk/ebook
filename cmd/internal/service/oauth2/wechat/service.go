package wechat

import "net/http"

type Service struct {
	appId     string
	appSecret string
	client    *http.Client
	//cmd       redis.Cmdable
}
