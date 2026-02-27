package game

import (
	"log"
	"sync"
)

// PacketHandler 数据包处理函数
type PacketHandler func(pkt *Packet)

// Dispatcher 数据包分发器
// 注册不同 PacketID 的处理函数, 从 ReadLoop 中分发
type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[uint32]PacketHandler
	fallback PacketHandler
}

// NewDispatcher 创建分发器
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		handlers: make(map[uint32]PacketHandler),
	}
}

// Register 注册特定包ID的处理函数
func (d *Dispatcher) Register(packetID uint32, handler PacketHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[packetID] = handler
}

// SetFallback 设置默认处理函数 (未注册包ID时调用)
func (d *Dispatcher) SetFallback(handler PacketHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.fallback = handler
}

// Dispatch 分发数据包到对应的处理函数
func (d *Dispatcher) Dispatch(pkt *Packet) {
	d.mu.RLock()
	handler, ok := d.handlers[pkt.ID]
	fallback := d.fallback
	d.mu.RUnlock()

	if ok {
		handler(pkt)
	} else if fallback != nil {
		fallback(pkt)
	} else {
		log.Printf("[game] 未处理的数据包: %s len=%d",
			PacketName(pkt.ID), len(pkt.Body))
	}
}

// HandleRaw 处理原始 WebSocket 二进制消息
// 自动解析并分发所有包
func (d *Dispatcher) HandleRaw(data []byte) {
	packets := ParsePackets(data)
	for _, pkt := range packets {
		d.Dispatch(pkt)
	}
}

// ConnWithDispatcher 使用分发器的游戏连接封装
// 在 Conn 基础上增加协议层的包解析和分发
func (c *Conn) ReadLoopDispatch(dispatcher *Dispatcher) {
	c.ReadLoop(func(msgType int, data []byte) {
		dispatcher.HandleRaw(data)
	})
}

// SendPacket 发送游戏数据包
func (c *Conn) SendPacket(pkt *Packet) error {
	return c.Send(pkt.Encode())
}

// SendHeartbeat 发送心跳包
func (c *Conn) SendHeartbeat() error {
	return c.SendPacket(&Packet{ID: PktHeartbeat, Body: nil})
}
