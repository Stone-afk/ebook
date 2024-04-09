package simpleim

import (
	"context"
	"ebook/cmd/pkg/logger"
	"ebook/cmd/pkg/saramax"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/ecodeclub/ekit/syncx"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

// Conn 稍微做一个封装
type Conn struct {
	*websocket.Conn
}

func (c *Conn) Send(msg Message) error {
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.WriteMessage(websocket.TextMessage, val)
}

type WsGateway struct {
	l logger.Logger
	// 连接了这个实例的客户端
	// 这里我们用 uid 作为 key
	// 实践中要考虑到不同的设备，
	// 那么这个 key 可能是一个复合结构，例如 uid + 设备
	conns *syncx.Map[int64, *Conn]
	svc   *IMService

	client     sarama.Client
	instanceId string
	upgrader   *websocket.Upgrader
}

// Start 在这个启动的时候，监听 websocket 的请求，然后转发到后端
func (g *WsGateway) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", g.wsHandler)
	err := g.subscribeMsg()
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, mux)
}

func (g *WsGateway) wsHandler(writer http.ResponseWriter, request *http.Request) {
	panic("")
}

// Uid 一般是从 jwt token 或者 session 里面取出来
// 这里模拟从 header 里面读取出来
func (g *WsGateway) Uid(req *http.Request) int64 {

	// 拿到 token
	//token := strings.TrimLeft(req.Header.Get("Authorization"), "Bearer ")
	// jwt 解析
	// jwt.Parse
	// req.Cookie("sess_id")

	uidStr := req.Header.Get("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)
	return uid
}

func (g *WsGateway) subscribeMsg() error {
	// 用 instance id 作为消费者组
	// 不像业务里面，同样的节点同一个消费者组
	// 每个节点单独的消费者组
	cg, err := sarama.NewConsumerGroupFromClient(g.instanceId, g.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{eventName},
			saramax.NewHandler[Event](g.l, g.consume))
		if err != nil {
			log.Println("退出监听消息循环", err)
		}
	}()
	return nil
}

func (g *WsGateway) consume(msg *sarama.ConsumerMessage, evt Event) error {
	// 转发
	// 我怎么知道，这个 receiver 有没有连上我？
	// 多端同步的时候，还需要知道哪个设备连上了我
	receiverConn, ok := g.conns.Load(evt.Receiver)
	if !ok {
		return nil
	}
	return receiverConn.Send(evt.Msg)
}
