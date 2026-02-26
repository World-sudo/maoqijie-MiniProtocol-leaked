package register

// NativeCallback 回调类型常量
// 逆向自 libiworld.dll: NativeCalledLoginManager 的回调消息类型
const (
	// CallbackFeatureResult 功能结果回调
	CallbackFeatureResult = "NativeCallFeatureResult"
	// CallbackViewOp 视图操作回调
	CallbackViewOp = "NativeCalledViewOp"
	// CallbackPermission 权限回调
	CallbackPermission = "NativePermissionCallback"
	// CallbackDomainLoginView 域名登录视图
	CallbackDomainLoginView = "DomainLoginView"
	// CallbackDomainLoginResult 域名登录结果
	CallbackDomainLoginResult = "DomainLoginResult"
)

// SDKLoginResult SDK登录回调结果
// 逆向自 libiworld.dll: SDKLoginCallBack(%d, [[%s]], [[%s]], %d)
type SDKLoginResult struct {
	Status int    // 状态码
	UID    string // 用户ID
	Token  string // 登录令牌 (JWT)
	Type   int    // 登录类型
}

// PlatformLoginResult 平台SDK登录结果
// 逆向自 libiworld.dll: login Result:%d uid:%s token:%s username:%s
type PlatformLoginResult struct {
	Result   int
	UID      string
	Token    string
	Username string
}

// LoginRoute H5登录页面路由类型
// 逆向自 LJ#45 uiJump/BasicSupportUI
const (
	RouteRegister    = "reg"
	RoutePassword    = "pwd"
	RouteSetPassword = "_setpassword"
	RouteCaptcha     = "captcha"
	RouteOneClick    = "oneclick"
	RouteJuhePhone   = "juhe_phone"
	RouteSecurity    = "_security"
	Route4399        = "4399_f"
	RouteQQWeChat    = "qq_wechat"
)

// LoginAuthMode 认证模式
// 逆向自 pcapng JWT: auth 字段
const (
	AuthModeWeb = "web"
)
