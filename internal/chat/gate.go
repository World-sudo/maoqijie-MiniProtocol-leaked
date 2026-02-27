package chat

import (
	"fmt"
	"log"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// Gate 聊天WebSocket长连接
type Gate struct {
	cred   *auth.Credential
	conn   *websocket.Conn
	host   string
	stopCh chan struct{}
}

// NewGate 创建聊天网关
func NewGate(cred *auth.Credential, host string) *Gate {
	return &Gate{cred: cred, host: host, stopCh: make(chan struct{})}
}

// Connect 建立WebSocket长连接
func (g *Gate) Connect() error {
	ts := time.Now().Unix()
	authSig := g.cred.ChatAuthAt(ts)

	params := url.Values{}
	params.Set("uid", strconv.FormatInt(g.cred.Uin, 10))
	params.Set("token", g.cred.IMJWT)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("auth", authSig)
	params.Set("cltversion", strconv.Itoa(config.CltVersion))
	params.Set("apiid", strconv.Itoa(config.APIID))
	params.Set("reconnect", "0")

	u := fmt.Sprintf("ws://%s%s?%s", g.host, config.ChatGatePath, params.Encode())

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(u, nil)
	if err != nil {
		return fmt.Errorf("连接聊天网关失败: %w", err)
	}
	g.conn = conn
	log.Printf("[chat] 已连接聊天网关: %s", g.host)
	go g.heartbeatLoop()
	return nil
}

// ReadLoop 读取消息循环
func (g *Gate) ReadLoop(handler func(msgType int, data []byte)) {
	if g.conn == nil {
		return
	}
	for {
		msgType, data, err := g.conn.ReadMessage()
		if err != nil {
			log.Printf("[chat] 读取消息失败: %v", err)
			return
		}
		handler(msgType, data)
	}
}

// Send 发送消息
func (g *Gate) Send(data []byte) error {
	if g.conn == nil {
		return fmt.Errorf("聊天网关未连接")
	}
	return g.conn.WriteMessage(websocket.BinaryMessage, data)
}

// Close 关闭连接
func (g *Gate) Close() error {
	select {
	case <-g.stopCh:
	default:
		close(g.stopCh)
	}
	if g.conn == nil {
		return nil
	}
	return g.conn.Close()
}

// heartbeatLoop 定期发送心跳帧保持连接
func (g *Gate) heartbeatLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-g.stopCh:
			return
		case <-ticker.C:
			frame := NewHeartbeatFrame()
			if err := g.Send(frame.Encode()); err != nil {
				log.Printf("[chat] 心跳发送失败: %v", err)
				return
			}
			log.Println("[chat] 心跳已发送")
		}
	}
}
