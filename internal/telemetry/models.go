package telemetry

// Event 遥测事件结构体
// 逆向自 pcapng: POST /miniworld 的 json3 字段
type Event struct {
	Age              int    `json:"age"`
	APIID            int    `json:"apiid"`
	Birthday         any    `json:"birthday"`
	ChannelAge       int    `json:"channel_age"`
	ChannelRealname  int    `json:"channel_isrealname"`
	CltAPIID         int    `json:"cltapiid"`
	CltVersion       int    `json:"cltversion"`
	Country          string `json:"country"`
	CountryAuthPI    any    `json:"countryauthpi"`
	DeviceID         string `json:"device_id"`
	Env              int    `json:"env"`
	ID               int    `json:"id"`
	IP               string `json:"ip"`
	IsAdult          int    `json:"is_adult"`
	IsRealnamePass   int    `json:"is_realname_pass"`
	IsRealnameConf   int    `json:"is_realname_conf,omitempty"`
	Language         int    `json:"language"`
	Log              string `json:"log"`
	TS               int64  `json:"ts"`
	Uin              int64  `json:"uin"`
	UinRegTime       int64  `json:"uin_reg_time,omitempty"`
}

// 事件ID常量
const (
	EventLoginCheckConnect = 66000
)

// 事件日志名常量
// 逆向自 pcapng 遥测事件流
const (
	LogLoginCheckConnect = "login_check_connect"
	LogLoginCheck        = "login_check"
	LogLoginRequest      = "login_request"
	LogLoginSuccess      = "login_success"
	LogLoginFailed       = "login_failed"
)

// UIEvent UI埋点事件
// 逆向自 pcapng: 页面/组件级别的用户行为追踪
type UIEvent struct {
	Page      string `json:"page"`
	Component string `json:"component"`
	Action    string `json:"action"`
	TS        int64  `json:"time_stamp"`
}

// UI页面常量
const (
	PageStartLoading = "START_PAGE_LOADING"
	PageLogin        = "LOGIN"
	PageSignIn       = "SIGNIN_ACCOUNT_1"
	PageSignSuccess  = "SIGNIN_SUCCESS"
	PageLandingPage  = "MINI_LANDINGPAGE_START_1"
)

// UI组件常量
const (
	CompLoginButton = "LogInButton"
	CompSignButton  = "SigninButton"
	CompPwdInput    = "PsdPut"
	CompEnterGame   = "EnterGame"
	CompLoginCard   = "LoginCard"
	CompNotice      = "Notice"
)
