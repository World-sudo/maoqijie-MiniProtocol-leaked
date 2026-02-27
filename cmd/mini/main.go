package main

import (
	"flag"
	"fmt"
	"log"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/chat"
	"miniprotocol/internal/credit"
	"miniprotocol/internal/device"
	"miniprotocol/internal/game"
	"miniprotocol/internal/httpc"
	"miniprotocol/internal/moderation"
	"miniprotocol/internal/room"
	"miniprotocol/internal/telemetry"
	"miniprotocol/internal/version"
	"os"
	"os/signal"
)

func main() {
	uin := flag.Int64("uin", 0, "用户uin (迷你号)")
	password := flag.String("pwd", "", "登录密码")
	nativeLogin := flag.Bool("native", false, "使用原生登录 (wskacchm)")
	doRegister := flag.Bool("register", false, "注册新账号")
	deviceID := flag.String("device", "", "设备指纹 (WINxxxx)")
	jwt := flag.String("jwt", "", "登录JWT令牌")
	skipTelemetry := flag.Bool("skip-telemetry", false, "跳过遥测上报")
	skipChat := flag.Bool("skip-chat", false, "跳过聊天服务")
	skipGame := flag.Bool("skip-game", false, "跳过游戏连接")
	loginOnly := flag.Bool("login-only", false, "仅登录，不连接游戏服务")
	checkVersion := flag.Bool("version", false, "查询版本信息")
	checkText := flag.String("check-text", "", "内容审核文本")
	queryCredit := flag.Bool("credit", false, "查询信用分")
	flag.Parse()

	// 设备指纹
	if *deviceID == "" {
		*deviceID = device.Generate()
		log.Printf("[main] 生成设备指纹: %s", *deviceID)
	} else if !device.Validate(*deviceID) {
		log.Fatalf("[main] 设备指纹格式错误: %s", *deviceID)
	}

	client := httpc.New()

	// 版本查询模式
	if *checkVersion {
		runVersionCheck(client)
		return
	}

	// 注册模式
	if *doRegister {
		if *password == "" {
			fmt.Fprintln(os.Stderr, "注册需要密码: mini -register -pwd <密码>")
			os.Exit(1)
		}
		runRegister(client, *password, *deviceID)
		return
	}

	// 信用分查询 (需要uin)
	if *queryCredit {
		if *uin == 0 {
			fmt.Fprintln(os.Stderr, "查询信用分需要uin: mini -credit -uin <迷你号>")
			os.Exit(1)
		}
		runCreditQuery(client, *uin)
		return
	}

	// 登录模式需要 uin
	if *uin == 0 && *password != "" && !*nativeLogin {
		fmt.Fprintln(os.Stderr, "SSO登录需要uin: mini -uin <迷你号> -pwd <密码>")
		os.Exit(1)
	}

	if *uin == 0 && *jwt == "" && *password == "" {
		fmt.Fprintln(os.Stderr, "用法:")
		fmt.Fprintln(os.Stderr, "  注册:   mini -register -pwd <密码>")
		fmt.Fprintln(os.Stderr, "  登录:   mini -uin <迷你号> -pwd <密码> -native")
		fmt.Fprintln(os.Stderr, "  SSO:    mini -uin <迷你号> -pwd <密码>")
		fmt.Fprintln(os.Stderr, "  连接:   mini -uin <迷你号> -jwt <token>")
		fmt.Fprintln(os.Stderr, "  版本:   mini -version")
		fmt.Fprintln(os.Stderr, "  审核:   mini -uin <迷你号> -jwt <token> -check-text <文本>")
		fmt.Fprintln(os.Stderr, "  信用分: mini -credit -uin <迷你号>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// 登录模式
	if *password != "" {
		if *nativeLogin {
			token, err := runNativeLogin(client, *uin, *password, *deviceID)
			if err != nil {
				log.Fatalf("[login] 原生登录失败: %v", err)
			}
			log.Printf("[login] 原生登录成功! token: %s", token)
			*jwt = token
		} else {
			token, err := runSSOLogin(client, *uin, *password)
			if err != nil {
				log.Fatalf("[login] SSO登录失败: %v", err)
			}
			log.Printf("[login] SSO登录成功! JWT: %s", token)
			*jwt = token
		}
		if *loginOnly {
			return
		}
	}

	if *uin == 0 {
		fmt.Fprintln(os.Stderr, "需要uin才能连接游戏服务")
		os.Exit(1)
	}

	// 凭证
	cred := auth.NewCredential(*uin, *deviceID)
	if *jwt != "" {
		cred.SetLoginJWT(*jwt)
	}
	log.Printf("[main] 凭证: %s", cred)

	// 内容审核 (需要登录后)
	if *checkText != "" {
		runTextCheck(client, cred, *checkText)
		return
	}

	// 遥测上报
	if !*skipTelemetry {
		runTelemetry(client, cred)
	}

	// 房间配置
	runRoom(client, cred)

	// 聊天服务
	var gate *chat.Gate
	if !*skipChat {
		gate = runChat(client, cred)
	}

	// 游戏连接
	var gameConn *game.Conn
	if !*skipGame {
		gameConn = runGame()
	}

	// 等待中断信号
	log.Println("[main] 运行中，按 Ctrl+C 退出")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	log.Println("[main] 正在关闭...")
	if gate != nil {
		gate.Close()
	}
	if gameConn != nil {
		gameConn.Close()
	}
}

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
		log.Printf("[chat] 收到消息: type=%d len=%d", msgType, len(data))
	})

	return gate
}

func runGame() *game.Conn {
	conn, err := game.Connect()
	if err != nil {
		log.Printf("[game] 连接失败: %v", err)
		return nil
	}

	go conn.ReadLoop(func(msgType int, data []byte) {
		log.Printf("[game] 收到消息: type=%d len=%d", msgType, len(data))
	})

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
