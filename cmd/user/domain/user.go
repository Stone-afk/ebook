package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Nickname string
	Password string
	Phone    string
	AboutMe  string

	//OpenID  string
	//UnionID string

	Ctime    time.Time
	Birthday time.Time

	// 不要组合，万一你将来可能还有 DingDingInfo，里面有同名字段 UnionID
	WechatInfo WechatInfo
}

// WechatInfo 微信的授权信息
type WechatInfo struct {
	// OpenId 是应用内唯一
	OpenId string
	// UnionId 是整个公司账号内唯一
	UnionId string
}
