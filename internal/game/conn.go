package game

import (
	"encoding/json"
	"fmt"
	"log"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Conn 游戏逻辑WebSocket连接
type Conn struct {
	conn *websocket.Conn
}

// Connect 连接游戏逻辑服务器
func Connect() (*Conn, error) {
	u := fmt.Sprintf("ws://%s:%d/", config.GameHost, config.GamePort)
	origin := fmt.Sprintf("http://%s:%d", config.GameHost, config.GamePort)

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
		Subprotocols:     []string{config.GameSubProtocol},
	}

	header := http.Header{}
	header.Set("Origin", origin)

	conn, _, err := dialer.Dial(u, header)
	if err != nil {
		return nil, fmt.Errorf("连接游戏服务器失败: %w", err)
	}

	log.Printf("[game] 已连接: %s", u)
	return &Conn{conn: conn}, nil
}

// ReadLoop 读取消息循环
func (c *Conn) ReadLoop(handler func(msgType int, data []byte)) {
	for {
		msgType, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[game] 读取失败: %v", err)
			return
		}
		handler(msgType, data)
	}
}

// Send 发送消息
func (c *Conn) Send(data []byte) error {
	return c.conn.WriteMessage(websocket.BinaryMessage, data)
}

// ConnectWithAuth 连接游戏服务器并发送登录认证包
func ConnectWithAuth(cred *auth.Credential) (*Conn, error) {
	c, err := Connect()
	if err != nil {
		return nil, err
	}

	// 构造登录请求包 (body = JWT 认证信息)
	loginBody, _ := json.Marshal(map[string]any{
		"auth": cred.AuthString(),
		"uin":  cred.Uin,
	})
	loginPkt := &Packet{ID: PktLoginReq, Body: loginBody}
	if err := c.SendPacket(loginPkt); err != nil {
		c.Close()
		return nil, fmt.Errorf("发送登录包失败: %w", err)
	}
	log.Printf("[game] 已发送登录认证包")
	return c, nil
}

// Close 关闭连接
func (c *Conn) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}
