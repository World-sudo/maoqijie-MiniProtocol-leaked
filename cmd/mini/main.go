package main

import (
	"flag"
	"fmt"
	"log"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/chat"
	"miniprotocol/internal/device"
	"miniprotocol/internal/game"
	"miniprotocol/internal/httpc"
	"os"
	"os/signal"
)

func main() {
	// 基础
	uin := flag.Int64("uin", 0, "用户uin (迷你号)")
	password := flag.String("pwd", "", "登录密码")
	nativeLogin := flag.Bool("native", false, "使用原生登录")
	doRegister := flag.Bool("register", false, "注册新账号")
	deviceID := flag.String("device", "", "设备指纹 (WINxxxx)")
	jwt := flag.String("jwt", "", "登录JWT令牌")
	skipTelemetry := flag.Bool("skip-telemetry", false, "跳过遥测上报")
	skipChat := flag.Bool("skip-chat", false, "跳过聊天服务")
	skipGame := flag.Bool("skip-game", false, "跳过游戏连接")
	loginOnly := flag.Bool("login-only", false, "仅登录")
	checkVersion := flag.Bool("version", false, "查询版本信息")
	checkText := flag.String("check-text", "", "内容审核文本")
	queryCredit := flag.Bool("credit", false, "查询信用分")
	checkUpdate := flag.Bool("update", false, "检查热更新")

	// 好友
	showFriends := flag.Bool("friends", false, "查询好友列表")
	addFriend := flag.Int64("add-friend", 0, "添加好友 (uin)")
	delFriend := flag.Int64("del-friend", 0, "删除好友 (uin)")
	acceptFriend := flag.Int64("accept-friend", 0, "接受好友申请 (uin)")
	friendRequests := flag.Bool("friend-requests", false, "查看好友申请列表")
	onlineFriends := flag.Bool("online-friends", false, "查看在线好友")
	searchUser := flag.String("search", "", "搜索用户")

	// 资料
	queryProfile := flag.Int64("profile", 0, "查询用户资料 (uin)")
	queryCard := flag.Int64("card", 0, "查询用户名片 (uin)")
	setNick := flag.String("set-nick", "", "修改昵称")
	setSignFlag := flag.String("set-sign", "", "修改签名")
	setGender := flag.Int("set-gender", -1, "修改性别 (0未知/1男/2女)")
	setSkin := flag.Int("set-skin", -1, "修改皮肤ID")

	// 社交
	followUser := flag.Int64("follow", 0, "关注用户 (uin)")
	unfollowUser := flag.Int64("unfollow", 0, "取消关注 (uin)")
	likeUser := flag.Int64("like", 0, "点赞用户 (uin)")
	showFans := flag.Int64("fans", 0, "查询粉丝 (uin)")
	showFollowing := flag.Int64("following", 0, "查询关注列表 (uin)")

	// 邮件
	showMail := flag.Bool("mail", false, "查询邮件列表")
	mailRecvAll := flag.Bool("mail-recv", false, "领取所有邮件附件")
	delMail := flag.String("del-mail", "", "删除邮件 (ID)")

	// 排行榜
	showRank := flag.Bool("rank", false, "全服排行榜")
	rankFriend := flag.Bool("rank-friend", false, "好友排行榜")
	rankMap := flag.Bool("rank-map", false, "地图排行榜")
	rankWeekly := flag.Bool("rank-weekly", false, "周榜")

	// 防沉迷
	antiAddict := flag.Bool("antiaddict", false, "防沉迷查询")

	flag.Parse()

	// 设备指纹
	if *deviceID == "" {
		*deviceID = device.Generate()
		log.Printf("[main] 生成设备指纹: %s", *deviceID)
	} else if !device.Validate(*deviceID) {
		log.Fatalf("[main] 设备指纹格式错误: %s", *deviceID)
	}

	client := httpc.New()

	// 无需凭证的命令
	if *checkVersion {
		runVersionCheck(client)
		return
	}
	if *doRegister {
		if *password == "" {
			fmt.Fprintln(os.Stderr, "注册需要密码: mini -register -pwd <密码>")
			os.Exit(1)
		}
		runRegister(client, *password, *deviceID)
		return
	}
	if *queryCredit {
		requireUin(*uin, "信用分")
		runCreditQuery(client, *uin)
		return
	}
	if *checkUpdate {
		runUpdateCheck(client, *uin)
		return
	}
	if *antiAddict {
		requireUin(*uin, "防沉迷")
		runAntiAddict(client, *uin)
		return
	}

	// 需要登录
	if *uin == 0 && *password != "" && !*nativeLogin {
		fmt.Fprintln(os.Stderr, "SSO登录需要uin: mini -uin <迷你号> -pwd <密码>")
		os.Exit(1)
	}
	if *uin == 0 && *jwt == "" && *password == "" {
		printUsage()
		os.Exit(1)
	}

	// 登录
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

	cred := auth.NewCredential(*uin, *deviceID)
	if *jwt != "" {
		cred.SetLoginJWT(*jwt)
	}
	log.Printf("[main] 凭证: %s", cred)

	// 需要凭证的一次性命令
	if dispatchOneShot(client, cred, oneShotArgs{
		checkText: *checkText, showFriends: *showFriends, addFriend: *addFriend,
		delFriend: *delFriend, acceptFriend: *acceptFriend, friendRequests: *friendRequests,
		onlineFriends: *onlineFriends, searchUser: *searchUser,
		queryProfile: *queryProfile, queryCard: *queryCard, setNick: *setNick,
		setSign: *setSignFlag, setGender: *setGender, setSkin: *setSkin,
		followUser: *followUser, unfollowUser: *unfollowUser, likeUser: *likeUser,
		showFans: *showFans, showFollowing: *showFollowing,
		showMail: *showMail, mailRecvAll: *mailRecvAll, delMail: *delMail,
		showRank: *showRank, rankFriend: *rankFriend, rankMap: *rankMap, rankWeekly: *rankWeekly,
	}) {
		return
	}

	// 完整连接模式
	if !*skipTelemetry {
		runTelemetry(client, cred)
	}
	runRoom(client, cred)

	var gate *chat.Gate
	if !*skipChat {
		gate = runChat(client, cred)
	}
	var gameConn *game.Conn
	if !*skipGame {
		gameConn = runGame(cred)
	}

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

type oneShotArgs struct {
	checkText      string
	showFriends    bool
	addFriend      int64
	delFriend      int64
	acceptFriend   int64
	friendRequests bool
	onlineFriends  bool
	searchUser     string
	queryProfile   int64
	queryCard      int64
	setNick        string
	setSign        string
	setGender      int
	setSkin        int
	followUser     int64
	unfollowUser   int64
	likeUser       int64
	showFans       int64
	showFollowing  int64
	showMail       bool
	mailRecvAll    bool
	delMail        string
	showRank       bool
	rankFriend     bool
	rankMap        bool
	rankWeekly     bool
}

// dispatchOneShot 分发一次性命令，返回 true 表示已处理
func dispatchOneShot(client *httpc.Client, cred *auth.Credential, a oneShotArgs) bool {
	switch {
	case a.checkText != "":
		runTextCheck(client, cred, a.checkText)
	case a.showFriends:
		runFriendList(client, cred)
	case a.addFriend > 0:
		runFriendAdd(client, cred, a.addFriend)
	case a.delFriend > 0:
		runFriendDelete(client, cred, a.delFriend)
	case a.acceptFriend > 0:
		runFriendAccept(client, cred, a.acceptFriend)
	case a.friendRequests:
		runFriendRequests(client, cred)
	case a.onlineFriends:
		runOnlineFriends(client, cred)
	case a.searchUser != "":
		runFriendSearch(client, cred, a.searchUser)
	case a.queryProfile > 0:
		runProfileQuery(client, cred, a.queryProfile)
	case a.queryCard > 0:
		runProfileCard(client, cred, a.queryCard)
	case a.setNick != "":
		runSetNick(client, cred, a.setNick)
	case a.setSign != "":
		runSetSign(client, cred, a.setSign)
	case a.setGender >= 0:
		runSetGender(client, cred, a.setGender)
	case a.setSkin >= 0:
		runSetSkin(client, cred, a.setSkin)
	case a.followUser > 0:
		runFollow(client, cred, a.followUser)
	case a.unfollowUser > 0:
		runUnfollow(client, cred, a.unfollowUser)
	case a.likeUser > 0:
		runLike(client, cred, a.likeUser)
	case a.showFans > 0:
		runFansList(client, cred, a.showFans)
	case a.showFollowing > 0:
		runFollowingList(client, cred, a.showFollowing)
	case a.showMail:
		runMailList(client, cred)
	case a.mailRecvAll:
		runMailReceiveAll(client, cred)
	case a.delMail != "":
		runMailDelete(client, cred, a.delMail)
	case a.showRank:
		runRankGlobal(client, cred)
	case a.rankFriend:
		runRankFriend(client, cred)
	case a.rankMap:
		runRankMap(client, cred)
	case a.rankWeekly:
		runRankWeekly(client, cred)
	default:
		return false
	}
	return true
}

func requireUin(uin int64, label string) {
	if uin == 0 {
		fmt.Fprintf(os.Stderr, "%s需要uin: mini -uin <迷你号>\n", label)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "用法:")
	fmt.Fprintln(os.Stderr, "  注册:     mini -register -pwd <密码>")
	fmt.Fprintln(os.Stderr, "  原生登录: mini -uin <迷你号> -pwd <密码> -native")
	fmt.Fprintln(os.Stderr, "  SSO登录:  mini -uin <迷你号> -pwd <密码>")
	fmt.Fprintln(os.Stderr, "  完整连接: mini -uin <迷你号> -jwt <token>")
	fmt.Fprintln(os.Stderr, "  版本:     mini -version")
	fmt.Fprintln(os.Stderr, "  信用分:   mini -credit -uin <迷你号>")
	fmt.Fprintln(os.Stderr, "  防沉迷:   mini -antiaddict -uin <迷你号>")
	flag.PrintDefaults()
}
