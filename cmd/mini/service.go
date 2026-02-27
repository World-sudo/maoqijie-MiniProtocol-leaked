package main

import (
	"log"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/chat"
	"miniprotocol/internal/credit"
	"miniprotocol/internal/game"
	"miniprotocol/internal/httpc"
	"miniprotocol/internal/moderation"
	"miniprotocol/internal/room"
	"miniprotocol/internal/telemetry"
	"miniprotocol/internal/version"
)

func runTelemetry(client *httpc.Client, cred *auth.Credential) {
	reporter := telemetry.NewReporter(client, cred)
	evt := reporter.LoginCheckConnectEvent()
	if err := reporter.Report([]telemetry.Event{evt}); err != nil {
		log.Printf("[telemetry] 上报失败: %v", err)
	} else {
		log.Println("[telemetry] 登录检查事件已上报")
	}
}

func runRoom(client *httpc.Client, cred *auth.Credential) {
	svc := room.NewService(client, cred)
	cfg, err := svc.GetConfig()
	if err != nil {
		log.Printf("[room] 获取房间配置失败: %v", err)
		return
	}
	log.Printf("[room] 房间: %s (%s:%d), 代理: %s:%d",
		cfg.Config.RoomName,
		cfg.Config.Room.IP, cfg.Config.Room.Port,
		cfg.Config.Proxy.IP, cfg.Config.Proxy.Port)
}

func runChat(client *httpc.Client, cred *auth.Credential) *chat.Gate {
	allocator := chat.NewAllocator(client, cred)
	allocResp, err := allocator.Alloc()
	if err != nil {
		log.Printf("[chat] 分配节点失败: %v", err)
		return nil
	}
	log.Printf("[chat] 分配节点: %s", allocResp.Data.Host)

	gate := chat.NewGate(cred, allocResp.Data.Host)
	if err := gate.Connect(); err != nil {
		log.Printf("[chat] 连接网关失败: %v", err)
		return nil
	}

	go gate.ReadLoop(func(msgType int, data []byte) {
		frames := chat.ParseFrames(data)
		for _, f := range frames {
			switch f.Opcode {
			case chat.OpTextMsg:
				if msg, err := f.ParseTextMessage(); err == nil {
					log.Printf("[chat] 文本消息: %d→%d %s", msg.FromUin, msg.ToUin, msg.Content)
				}
			case chat.OpSystemNotify:
				if n, err := f.ParseSystemNotify(); err == nil {
					log.Printf("[chat] 系统通知: %s - %s", n.Title, n.Content)
				}
			case chat.OpHeartbeatAck:
				log.Println("[chat] 心跳回复")
			default:
				log.Printf("[chat] %s len=%d", chat.OpcodeName(f.Opcode), len(f.Body))
			}
		}
	})

	return gate
}

func runGame(cred *auth.Credential) *game.Conn {
	conn, err := game.ConnectWithAuth(cred)
	if err != nil {
		log.Printf("[game] 连接失败: %v", err)
		return nil
	}

	disp := game.NewDispatcher()
	disp.Register(game.PktLoginResp, func(pkt *game.Packet) {
		log.Printf("[game] 登录响应: len=%d", len(pkt.Body))
	})
	disp.Register(game.PktHeartbeatAck, func(pkt *game.Packet) {
		log.Println("[game] 心跳回复")
	})
	disp.Register(game.PktPlayerJoin, func(pkt *game.Packet) {
		log.Printf("[game] 玩家加入: len=%d", len(pkt.Body))
	})
	disp.Register(game.PktPlayerLeave, func(pkt *game.Packet) {
		log.Printf("[game] 玩家离开: len=%d", len(pkt.Body))
	})
	disp.Register(game.PktChatRecv, func(pkt *game.Packet) {
		log.Printf("[game] 游戏聊天: %s", string(pkt.Body))
	})
	disp.SetFallback(func(pkt *game.Packet) {
		log.Printf("[game] %s len=%d", game.PacketName(pkt.ID), len(pkt.Body))
	})

	go conn.ReadLoopDispatch(disp)

	return conn
}

func runVersionCheck(client *httpc.Client) {
	checker := version.NewChecker(client)
	info, err := checker.GetVersionJSON()
	if err != nil {
		log.Printf("[version] 获取版本失败: %v", err)
		return
	}
	log.Printf("[version] 版本: %s (cltversion=%d)", info.Version, info.CltVersion)
	log.Printf("[version] 下载: %s", info.URL)
	if info.ForceUp != 0 {
		log.Printf("[version] 强制更新: %s", info.Desc)
	}
}

func runTextCheck(client *httpc.Client, cred *auth.Credential, text string) {
	checker := moderation.NewChecker(client, cred)
	result, err := checker.CheckText(text, "chat")
	if err != nil {
		log.Printf("[moderation] 审核失败: %v", err)
		return
	}
	log.Printf("[moderation] code=%d pass=%v reason=%s",
		result.Code, result.Data.Pass, result.Data.Reason)
}

func runCreditQuery(client *httpc.Client, uin int64) {
	c := credit.NewClient(client)
	result, err := c.QueryScore(uin)
	if err != nil {
		log.Printf("[credit] 查询失败: %v", err)
		return
	}
	log.Printf("[credit] code=%d score=%d limited=%v",
		result.Code, result.Data.Score, result.Data.Limited)
}
