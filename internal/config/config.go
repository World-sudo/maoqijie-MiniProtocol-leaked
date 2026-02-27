package config

const (
	// 客户端标识
	UserAgent  = "Rainbow/1.0 (Windows_RT; U; Linux 6.2; zh)"
	AppVersion = "1.53.2"
	CltVersion = 79106
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

	// SSO 登录 (新版, 逆向自 sso.mini1.cn H5页面)
	SSOHost      = "sso.mini1.cn"
	SSOLoginAPI  = "wapi.mini1.cn"
	SSOLoginPath = "/login-service/api/web/v2/login"
	SSOAuthPath  = "/login-service/api/web/v1/auth/add"

	// 旧版注册/登录 (HTTPS, 已迁移但仍部分可用)
	RegisterWebHost  = "mnweb.mini1.cn"
	RegisterH5Host   = "h5.mini1.cn"
	RegisterPath     = "/register/"
	TextPwdLoginPath = "/account/TextPwdLogin"
	DomainLoginPath  = "/miniw/ldap/auth"

	// 原生登录 (MicroMiniNew.exe, 逆向自二进制)
	NativeAuthPath = "/login/auth_security"
	NativeAuthIP   = "120.24.63.165"
	NativeAuthPort = 14000

	// RPC API (逆向自 LJ#7 rpc_do_http_post)
	RPCAPIPath   = "/api/v1"
	RPCTestHost1 = "124.71.98.30:8089"
	RPCTestHost2 = "116.205.254.139"

	// SMS/Email 验证 (逆向自 MicroMiniNew.exe)
	SMSSendPath    = "/sms/smssend/"
	SMSVerifyPath  = "/sms/smsverify/"
	EmailSendPath  = "/email/emailsend/"
	EmailVerifyPath = "/email/emailverify/"
	SMSCheckType   = "2"
	SMSID          = "461053"

	// 通道管理
	ChannelHost     = "wskacchm.mini1.cn"
	ChannelPortPre  = 14130
	ChannelPortPost = 14120

	// 运营配置
	OperateCDNHost = "operate2cdn.mini1.cn"

	// 账号服务器 (DNS缓存自libiworld.dll iworld.cfg)
	AccountServer = "account.svr.mini1.cn"

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

	// 签名类型 (逆向自 libiworld.dll HttpSignMgr)
	SignTypeMD5  = "md5"
	SignTypeSHA1 = "sha1"

	// Sentry DSN
	SentryDSN    = "https://98bca67bdc8939da32b4d77b923d40e4@miniwsentry.mini1.cn/2"
	SSOSentryDSN = "https://cc07236fb202454e8d835e906259228a@cloud-sentry.mini1.cn/41"

	// AES-128加密 (bdinfo参数加密, 逆向自 LJ#137)
	AESKey = "7q1WyNG3dE3CRy85"
	AESIV  = "Utz92Hjrky1XAX1B"

	// AES-256 auth_test 密钥 (逆向自 LJ#222)
	AuthTestKey = "6tZwR6zAcAkcj2NxMYOBuU1sCl8bphyH"

	// 原生登录响应解密 AES-256-CBC 密钥 (逆向自 MicroMiniNew.exe)
	// 初始化代码: 0x0064C440 -> push 0x65E104 -> "fcafc12e17b93a30a8998fcbc7d5c786"
	// 存储位置: g_aesKey 全局变量 0x6EB9DC
	NativeRespAESKey = "fcafc12e17b93a30a8998fcbc7d5c786"

	// 原生登录响应解密 AES-256-CBC IV (逆向自 MicroMiniNew.exe)
	// 初始化代码: 0x0064C460 -> push 0x65E128 -> "624df8d86de5dc35"
	// 存储位置: g_aesIV 全局变量 0x6EB9F8
	NativeRespAESIV = "624df8d86de5dc35"

	// 第三方OAuth
	WeChatAppID = "wx0344e7ba7bfcacaf"
	QQAppID     = "101901986"

	// DomainLogin常量
	DomainLoginHash = "f5711eb1640712de051e5aedc35329c3"

	// 固定认证MD5哈希 (逆向自 LJC字节码多处出现)
	AuthHashFixed = "763f86ba71a337e1681872f12e23b411" // &auth= 参数
	CTHashFixed   = "3dbc5f33add11d1af78ba2af365e095"  // &cthash= 参数
	ThirdHash     = "583006e8867d41f6b17e431d42c8b7e7" // 第三方登录
	BaseEngHash   = "7788ff50ea1eb307153b5202bc2c1477" // libMiniBaseEngine
	PkgHash       = "8aa7c844f44dca9a0af98edc49759b01" // game_script.pkg

	// RPC代理路径 (逆向自 patch_game_script.pkg)
	RPCProxyPath   = "/_proxy"
	SocialProxyCmd = "/social_proxy"

	// 社区/API服务
	ShequHost   = "shequ.mini1.cn"
	CreditAPI   = "credit-api.mini1.cn"
	MiniPalHost = "minipal.mini1.cn"

	// GeeTest V4 (逆向自 sso.mini1.cn JS)
	GeeTestCaptchaURL = "gcaptcha4.geetest.com"
	GeeTestCaptchaID  = "57157b87c9788ae72be45a2c79c6dd1c"

	// 极验 secret_id (逆向自 LJ#146, 按平台区分)
	GeeTestSecretAndroid = "7fd790bf9e3f8b31b9c477af7f569b74"
	GeeTestSecretApple   = "8faa3850bb007ddcefc08bcb0fcf34f2"

	// MicroMiniNew 签名salt (逆向自二进制 0x0065EEF8, 服务器不接受)
	NativeSignSalt = "c8c93222583741bd828579b3d3efd43b"

	// 原生认证服务器实际使用的签名salt (逆向自运行时 .data 段 0x006EBA14)
	// 签名格式: md5("source=mini_micro&target=<target>&time=<ts>" + salt)
	// JSON字段名为 "auth" (不是 "sign")
	NativeServerSalt = "2ddb7619717147439c83ab022e9d4d38"

	// 登录方式 (逆向自 MicroMiniNew.exe)
	LoginTypeTextPwd    = "TextPasswordLogin"
	LoginTypeDigitalPwd = "DigitalPasswordLogin"
	LoginTypeCreate     = "CreateAccount"

	// 认证方式 (逆向自 MicroMiniNew.exe)
	AuthModePasswd    = "passwd_auth"
	AuthModeAuthInfo  = "authinfo_auth"
	AuthModeQuestion  = "question_auth"
	AuthModePhoneEdu  = "phone_login_education"

	// CrashSight 崩溃上报
	CrashSightAppID = "0c2abe1373"
	CrashSightURL   = "pc.crashsight.qq.com"
)
