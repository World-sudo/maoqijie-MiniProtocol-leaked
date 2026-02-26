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

	// 第三方OAuth
	WeChatAppID = "wx0344e7ba7bfcacaf"
	QQAppID     = "101901986"

	// DomainLogin常量
	DomainLoginHash = "f5711eb1640712de051e5aedc35329c3"

	// GeeTest V4 (逆向自 sso.mini1.cn JS)
	GeeTestCaptchaURL = "gcaptcha4.geetest.com"
	GeeTestCaptchaID  = "57157b87c9788ae72be45a2c79c6dd1c"

	// MicroMiniNew 签名salt (逆向自二进制)
	NativeSignSalt = "c8c93222583741bd828579b3d3efd43b"

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
