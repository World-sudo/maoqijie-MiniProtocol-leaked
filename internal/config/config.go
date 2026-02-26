package config

const (
	// 客户端标识
	UserAgent  = "Rainbow/1.0 (Windows_RT; U; Linux 6.2; zh)"
	AppVersion = "1.53.1"
	CltVersion = 79105
	APIID      = 110

	// 游戏逻辑服务器
	GameHost = "cn-logic6.mini1.cn"
	GamePort = 4009

	// 遥测上报
	TelemetryHost    = "tj3.mini1.cn"
	TelemetryAltHost = "tj.mini1.cn"
	TelemetryPath    = "/miniworld"

	// 聊天服务
	ChatAllocHost = "chatpush.mini1.cn"
	ChatAllocPort = 19601
	ChatGatePort  = 19701
	ChatAllocPath = "/minilb/alloc"
	ChatGatePath  = "/minigate/gate"
	ChatRPCPath   = "/minilb/rpc"

	// 房间服务
	RoomHost = "openroom.mini1.cn"
	RoomPort = 8080
	RoomPath = "/server/room"

	// 注册/登录 (HTTPS)
	RegisterWebHost  = "mnweb.mini1.cn"
	RegisterH5Host   = "h5.mini1.cn"
	RegisterPath     = "/register/"
	TextPwdLoginPath = "/account/TextPwdLogin"
	DomainLoginPath  = "/miniw/ldap/auth"

	// 通道管理
	ChannelHost     = "wskacchm.mini1.cn"
	ChannelPortPre  = 14130
	ChannelPortPost = 14120

	// 运营配置
	OperateCDNHost = "operate2cdn.mini1.cn"

	// WebSocket子协议
	GameSubProtocol = "default-protocol"

	// 认证前缀
	AuthPrefix = "switchAccountByAuthInfo_reg###"

	// JWT签发来源
	JWTSource = "man_machine.login_v3"

	// 自定义HTTP头
	HeaderAuth    = "MN-AUTH"
	HeaderToken   = "MN-TOKEN"
	HeaderPayload = "MN-PAYLOAD"

	// Sentry DSN
	SentryDSN = "https://98bca67bdc8939da32b4d77b923d40e4@miniwsentry.mini1.cn/2"
)
