package register

// 服务端错误码 (字符串类型)
// 逆向自 LJ#138 错误码枚举
const (
	ErrNeedCaptcha    = "NEED_CAPTCHA"
	ErrAuthInfoNotOK  = "ACCOUNT_AUTHINFO_NOT_OK"
	ErrPasswdNotOK    = "ACCOUNT_PASSWD_NOT_OK"
	ErrDBNotFound     = "ACCOUNT_DB_NOT_FOUND"
	ErrRiskLevel      = "RISK_LEVEL"
	ErrLimit          = "LIMIT"
	ErrForMore        = "FOR_MORE"
	ErrBind           = "BIND"
	ErrAlready        = "ALREADY"
	ErrNoQuick        = "NO_QUICK"
	ErrToken          = "TOKEN"
	ErrSign           = "SIGN"
	ErrParams         = "PARAMS"
	ErrFreeze         = "FREEZE"
	ErrSafe           = "SAFE"
	ErrLuaTaskTimeout = "LUA_TASK_TIME_OUT"
)

// SSO 数字错误码 (逆向自 sso.mini1.cn JS + 抓包)
const (
	CodeSuccess         = 0    // 登录成功
	CodeNeedCaptcha     = 7211 // 需要GeeTest验证码
	CodeAccountCanceled = 4104 // 账号已申请注销
	CodeAccountBanned   = 4075 // 账号被封禁
)

// 原生认证数字错误码 (逆向自 wskacchm.mini1.cn:14130 响应)
const (
	CodeNativeAuthFailed     = 7012 // 认证失败 (密码错误)
	CodeNativeAPIIDInvalid   = 4001 // apiid 不正确
	CodeNativeServiceNotReady = 1001 // 服务未准备好 / 缺少注册信息
)

// SMS 预检查状态 (逆向自 LJ#136/LJ#148)
const (
	PrecheckSMS   = "PRECHECK_SMS"
	PrecheckAPIID = "PRECHECK_APIID_"
	PrecheckEmail = "PRECHECK_EMAIL_HAS_ALREADY_BEAN_B"
)
