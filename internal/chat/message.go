package chat

// 聊天消息协议解析
// 逆向推断自 libiworld.dll IM 模块:
//   消息帧格式: [2B opcode][2B bodyLen][body...]
//   WebSocket 二进制消息, 小端字节序
//
// 已知消息类型推断自 libiworld.dll 字符串:
//   heartbeat, text_msg, system_notify, friend_online, group_msg,
//   join_room, leave_room, kick_notify
//
// + pcapng 抓包观察到的 WebSocket 帧结构

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// 消息 Opcode 常量
// 逆向推断自 libiworld.dll IM 消息派发表
const (
	OpHeartbeat     uint16 = 0x0001 // 心跳 (客户端→服务器)
	OpHeartbeatAck  uint16 = 0x0002 // 心跳回复
	OpTextMsg       uint16 = 0x0010 // 文本消息
	OpImageMsg      uint16 = 0x0011 // 图片消息
	OpVoiceMsg      uint16 = 0x0012 // 语音消息
	OpSystemNotify  uint16 = 0x0020 // 系统通知
	OpFriendOnline  uint16 = 0x0030 // 好友上线
	OpFriendOffline uint16 = 0x0031 // 好友下线
	OpGroupMsg      uint16 = 0x0040 // 群组消息
	OpJoinRoom      uint16 = 0x0050 // 加入房间
	OpLeaveRoom     uint16 = 0x0051 // 离开房间
	OpKickNotify    uint16 = 0x0052 // 踢出通知
	OpReadReceipt   uint16 = 0x0060 // 已读回执
	OpTyping        uint16 = 0x0061 // 正在输入
)

// Frame 聊天消息帧
// 格式: [Opcode:2B LE][BodyLen:2B LE][Body:BodyLen bytes]
type Frame struct {
	Opcode  uint16
	Body    []byte
}

// HeaderSize 帧头固定大小
const HeaderSize = 4

// ParseFrame 从二进制数据解析一个消息帧
// 返回帧和消费的字节数, 数据不足时返回 nil
func ParseFrame(data []byte) (*Frame, int) {
	if len(data) < HeaderSize {
		return nil, 0
	}

	opcode := binary.LittleEndian.Uint16(data[0:2])
	bodyLen := binary.LittleEndian.Uint16(data[2:4])
	totalLen := HeaderSize + int(bodyLen)

	if len(data) < totalLen {
		return nil, 0
	}

	body := make([]byte, bodyLen)
	copy(body, data[HeaderSize:totalLen])

	return &Frame{Opcode: opcode, Body: body}, totalLen
}

// ParseFrames 从二进制数据解析所有完整帧
func ParseFrames(data []byte) []*Frame {
	var frames []*Frame
	offset := 0
	for offset < len(data) {
		frame, consumed := ParseFrame(data[offset:])
		if frame == nil {
			break
		}
		frames = append(frames, frame)
		offset += consumed
	}
	return frames
}

// Encode 将帧编码为二进制数据
func (f *Frame) Encode() []byte {
	buf := make([]byte, HeaderSize+len(f.Body))
	binary.LittleEndian.PutUint16(buf[0:2], f.Opcode)
	binary.LittleEndian.PutUint16(buf[2:4], uint16(len(f.Body)))
	copy(buf[HeaderSize:], f.Body)
	return buf
}

// TextMessage 文本消息体
type TextMessage struct {
	FromUin  int64  `json:"from_uin"`
	ToUin    int64  `json:"to_uin"`
	Content  string `json:"content"`
	MsgID    string `json:"msg_id"`
	SendTime int64  `json:"send_time"`
}

// SystemNotification 系统通知体
type SystemNotification struct {
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}

// FriendStatus 好友状态变更
type FriendStatus struct {
	Uin    int64 `json:"uin"`
	Online bool  `json:"online"`
}

// RoomEvent 房间事件
type RoomEvent struct {
	Uin      int64  `json:"uin"`
	NickName string `json:"nick_name"`
	RoomID   string `json:"room_id"`
}

// ParseTextMessage 解析文本消息体
func (f *Frame) ParseTextMessage() (*TextMessage, error) {
	if f.Opcode != OpTextMsg {
		return nil, fmt.Errorf("opcode不匹配: 期望0x%04X 实际0x%04X", OpTextMsg, f.Opcode)
	}
	var msg TextMessage
	if err := json.Unmarshal(f.Body, &msg); err != nil {
		return nil, fmt.Errorf("解析文本消息失败: %w", err)
	}
	return &msg, nil
}

// ParseSystemNotify 解析系统通知
func (f *Frame) ParseSystemNotify() (*SystemNotification, error) {
	if f.Opcode != OpSystemNotify {
		return nil, fmt.Errorf("opcode不匹配: 期望0x%04X 实际0x%04X", OpSystemNotify, f.Opcode)
	}
	var notify SystemNotification
	if err := json.Unmarshal(f.Body, &notify); err != nil {
		return nil, fmt.Errorf("解析系统通知失败: %w", err)
	}
	return &notify, nil
}

// NewTextFrame 构建文本消息帧
func NewTextFrame(toUin int64, content string) *Frame {
	msg := TextMessage{
		ToUin:   toUin,
		Content: content,
	}
	body, _ := json.Marshal(msg)
	return &Frame{Opcode: OpTextMsg, Body: body}
}

// NewHeartbeatFrame 构建心跳帧
func NewHeartbeatFrame() *Frame {
	return &Frame{Opcode: OpHeartbeat, Body: nil}
}

// OpcodeName 获取 opcode 的可读名称
func OpcodeName(op uint16) string {
	switch op {
	case OpHeartbeat:
		return "Heartbeat"
	case OpHeartbeatAck:
		return "HeartbeatAck"
	case OpTextMsg:
		return "TextMsg"
	case OpImageMsg:
		return "ImageMsg"
	case OpVoiceMsg:
		return "VoiceMsg"
	case OpSystemNotify:
		return "SystemNotify"
	case OpFriendOnline:
		return "FriendOnline"
	case OpFriendOffline:
		return "FriendOffline"
	case OpGroupMsg:
		return "GroupMsg"
	case OpJoinRoom:
		return "JoinRoom"
	case OpLeaveRoom:
		return "LeaveRoom"
	case OpKickNotify:
		return "KickNotify"
	case OpReadReceipt:
		return "ReadReceipt"
	case OpTyping:
		return "Typing"
	default:
		return fmt.Sprintf("Unknown(0x%04X)", op)
	}
}
