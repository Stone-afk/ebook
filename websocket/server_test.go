package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	//upgrader := &websocket.Upgrader{}
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {

	})
}

type Ws struct {
	*websocket.Conn
}

func (ws *Ws) WriteString(data string) error {
	err := ws.WriteMessage(websocket.TextMessage, []byte(data))
	return err
}

//
//func Read() {
//	var conn net.Conn
//
//	for {
//		// 如果每一次都创建这个 buffer，
//		buffer := pool.Get()
//		conn.Read(buffer)
//		// 用完了
//		pool.Put(buffer)
//	}
//}