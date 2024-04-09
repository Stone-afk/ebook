package simpleim

// MessageV1 和前端约定好，具体的消息的内容的格式
//type MessageV1 struct {
//
//	// 这个是前端的序列号
//	// 不要求全局唯一的，正常只要当下这个 websocket 唯一就可以
//	Seq string
//
//	// 谁发的？
//	// 能不能是前端传过来的？
//	// Sender int64
//
//	// 发给谁
//	// cid channel id(group id)，聊天 ID
//	// 单聊，也是用聊天 ID
//	Cid int64
//	// 内容
//	// Type 这个消息是什么消息
//	// 这个是你 IM 内部的类型
//	// type = "video", => content = url/资源标识符 key
//	// content 不可能是视频本身
//	// {"title": "GO从入门到入土", Addr: "https://oss.aliyun.com/im/resource/abc"}
// @某人 {"metions": []int64, "text": }
//	Type string
//	// 你有文本消息，你有图片消息，你有视频消息
//	// 你这个 Content 究竟是什么？
//	Content string
//
//	// 万一你每个消息都要校验 token，可以在这里带
//	//Token string
//}

type Message struct {
	// 发过来的消息的序列号
	// 用于前后端关联消息
	Seq string
	// 这个是后端的 ID
	// 前端有时候支持引用功能，转发功能的时候，会需要这个 ID
	ID int64
	// 用来标识不同的消息类型
	// 文本消息，视频消息
	// 系统消息（后端往前端发的，跟 IM 本身管理有关的消息）
	Type    string
	Content string
	// 聊天 ID，注意，正常来说这里不是记录目标用户 ID
	// 而是记录代表了这个聊天的 ID
	Cid int64
}

type Event struct {
	Msg Message
	// 接收者
	Receiver int64
	// 发送的 device
	Device string
}

// EventV1 扩散只会和你有多少接入节点有关
// 和群里面有多少人无关
// 注册与发现机制，那么你就可以精确控制，转发到哪些节点
type EventV1 struct {
	Msg       Message
	Receivers []int64
}

const eventName = "simple_im_msg"
