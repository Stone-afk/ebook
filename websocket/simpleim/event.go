package simpleim

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
