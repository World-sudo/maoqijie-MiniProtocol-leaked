package register

// 服务端错误码
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
