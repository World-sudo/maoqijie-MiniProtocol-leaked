package game

// 游戏协议包定义
// 逆向推断自 libiworld.dll + MicroMiniNew.exe 网络层:
//   包格式: [4B PacketID LE][4B BodyLen LE][Body...]
//   WebSocket 二进制帧, 小端字节序
//
// 包类型推断自 libiworld.dll 字符串表:
//   C2S: login_req, heartbeat, player_move, block_change, chat_send
//   S2C: login_resp, heartbeat_ack, player_sync, world_data, entity_spawn
//
// 观察: 游戏连接建立后首个包为 login_req (携带JWT认证)

import (
	"encoding/binary"
	"fmt"
)

// 包ID常量 (C2S = 客户端→服务器, S2C = 服务器→客户端)
// 逆向推断自 libiworld.dll 消息派发表
const (
	// 连接/认证
	PktLoginReq     uint32 = 0x0001 // C2S 登录认证
	PktLoginResp    uint32 = 0x0002 // S2C 登录响应
	PktHeartbeat    uint32 = 0x0003 // C2S 心跳
	PktHeartbeatAck uint32 = 0x0004 // S2C 心跳回复
	PktDisconnect   uint32 = 0x0005 // 双向 断开连接
	PktKick         uint32 = 0x0006 // S2C 踢出

	// 玩家
	PktPlayerMove   uint32 = 0x0010 // C2S 移动
	PktPlayerSync   uint32 = 0x0011 // S2C 其他玩家同步
	PktPlayerAction uint32 = 0x0012 // C2S 动作 (跳跃/攻击等)
	PktPlayerInfo   uint32 = 0x0013 // S2C 玩家信息
	PktPlayerJoin   uint32 = 0x0014 // S2C 玩家加入
	PktPlayerLeave  uint32 = 0x0015 // S2C 玩家离开

	// 世界
	PktWorldData    uint32 = 0x0020 // S2C 世界数据块
	PktBlockChange  uint32 = 0x0021 // C2S 方块变更
	PktBlockSync    uint32 = 0x0022 // S2C 方块同步
	PktChunkRequest uint32 = 0x0023 // C2S 请求区块数据
	PktChunkData    uint32 = 0x0024 // S2C 区块数据

	// 实体
	PktEntitySpawn   uint32 = 0x0030 // S2C 实体生成
	PktEntityDestroy uint32 = 0x0031 // S2C 实体销毁
	PktEntityMove    uint32 = 0x0032 // S2C 实体移动
	PktEntityAction  uint32 = 0x0033 // C2S 对实体操作

	// 聊天/UI
	PktChatSend     uint32 = 0x0040 // C2S 发送聊天
	PktChatRecv     uint32 = 0x0041 // S2C 接收聊天
	PktNotification uint32 = 0x0042 // S2C 通知消息

	// 物品/背包
	PktInventorySync  uint32 = 0x0050 // S2C 背包同步
	PktItemUse        uint32 = 0x0051 // C2S 使用物品
	PktItemDrop       uint32 = 0x0052 // C2S 丢弃物品
	PktItemPickup     uint32 = 0x0053 // C2S 拾取物品
	PktSlotChange     uint32 = 0x0054 // C2S 切换手持栏
)

// Packet 游戏数据包
// 格式: [PacketID:4B LE][BodyLen:4B LE][Body:BodyLen bytes]
type Packet struct {
	ID   uint32
	Body []byte
}

// PacketHeaderSize 包头固定大小
const PacketHeaderSize = 8

// ParsePacket 从二进制数据解析一个游戏数据包
// 返回包和消费的字节数, 数据不足时返回 nil
func ParsePacket(data []byte) (*Packet, int) {
	if len(data) < PacketHeaderSize {
		return nil, 0
	}

	id := binary.LittleEndian.Uint32(data[0:4])
	bodyLen := binary.LittleEndian.Uint32(data[4:8])
	totalLen := PacketHeaderSize + int(bodyLen)

	if len(data) < totalLen {
		return nil, 0
	}

	body := make([]byte, bodyLen)
	copy(body, data[PacketHeaderSize:totalLen])

	return &Packet{ID: id, Body: body}, totalLen
}

// ParsePackets 解析所有完整数据包
func ParsePackets(data []byte) []*Packet {
	var packets []*Packet
	offset := 0
	for offset < len(data) {
		pkt, consumed := ParsePacket(data[offset:])
		if pkt == nil {
			break
		}
		packets = append(packets, pkt)
		offset += consumed
	}
	return packets
}

// Encode 将数据包编码为二进制
func (p *Packet) Encode() []byte {
	buf := make([]byte, PacketHeaderSize+len(p.Body))
	binary.LittleEndian.PutUint32(buf[0:4], p.ID)
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(p.Body)))
	copy(buf[PacketHeaderSize:], p.Body)
	return buf
}

// PacketName 获取包ID的可读名称
func PacketName(id uint32) string {
	names := map[uint32]string{
		PktLoginReq: "LoginReq", PktLoginResp: "LoginResp",
		PktHeartbeat: "Heartbeat", PktHeartbeatAck: "HeartbeatAck",
		PktDisconnect: "Disconnect", PktKick: "Kick",
		PktPlayerMove: "PlayerMove", PktPlayerSync: "PlayerSync",
		PktPlayerAction: "PlayerAction", PktPlayerInfo: "PlayerInfo",
		PktPlayerJoin: "PlayerJoin", PktPlayerLeave: "PlayerLeave",
		PktWorldData: "WorldData", PktBlockChange: "BlockChange",
		PktBlockSync: "BlockSync", PktChunkRequest: "ChunkRequest",
		PktChunkData: "ChunkData",
		PktEntitySpawn: "EntitySpawn", PktEntityDestroy: "EntityDestroy",
		PktEntityMove: "EntityMove", PktEntityAction: "EntityAction",
		PktChatSend: "ChatSend", PktChatRecv: "ChatRecv",
		PktNotification: "Notification",
		PktInventorySync: "InventorySync", PktItemUse: "ItemUse",
		PktItemDrop: "ItemDrop", PktItemPickup: "ItemPickup",
		PktSlotChange: "SlotChange",
	}
	if name, ok := names[id]; ok {
		return name
	}
	return fmt.Sprintf("Unknown(0x%04X)", id)
}
