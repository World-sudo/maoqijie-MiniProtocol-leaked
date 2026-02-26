package telemetry

// Event 遥测事件结构体
type Event struct {
	Age              int    `json:"age"`
	APIID            int    `json:"apiid"`
	Birthday         int    `json:"birthday"`
	ChannelAge       int    `json:"channel_age"`
	ChannelRealname  int    `json:"channel_isrealname"`
	CltAPIID         int    `json:"cltapiid"`
	CltVersion       int    `json:"cltversion"`
	Country          string `json:"country"`
	CountryAuthPI    string `json:"countryauthpi"`
	DeviceID         string `json:"device_id"`
	Env              int    `json:"env"`
	ID               int    `json:"id"`
	IP               string `json:"ip"`
	IsAdult          int    `json:"is_adult"`
	IsRealnamePass   int    `json:"is_realname_pass"`
	Language         int    `json:"language"`
	Log              string `json:"log"`
	TS               int64  `json:"ts"`
	Uin              int64  `json:"uin"`
}

// 事件ID常量
const (
	EventLoginCheckConnect = 66000
)

// 事件日志名常量
const (
	LogLoginCheckConnect = "login_check_connect"
)
