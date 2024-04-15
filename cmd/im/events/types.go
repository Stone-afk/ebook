package events

type User struct {
	Id            int64
	Email         string
	Password      string
	Phone         string
	Birthday      string
	Nickname      string
	AboutMe       string
	WechatOpenId  string
	WechatUnionId string

	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
