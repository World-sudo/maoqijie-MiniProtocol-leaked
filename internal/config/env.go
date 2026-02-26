package config

// 环境配置 (逆向自 iworld.cfg game_env + libMiniBaseEngine.dll)
// game_env 控制服务器域名选择:
//   0  = 国内正式服
//   1  = 国内测试/预发布
//   10 = 海外正式服 (miniworldgame.com)
//   11 = 海外测试服

const (
	EnvDomestic    = 0
	EnvDomesticDev = 1
	EnvOversea     = 10
	EnvOverseaDev  = 11
)

// ServerSet 一组环境对应的服务器域名
type ServerSet struct {
	ShequHTTP  string // 社区/API HTTP (端口8080)
	ShequHTTPS string // 社区/API HTTPS (端口8081)
	WAPI       string // 核心 Web API
	Download   string // 资源下载
}

// Servers 各环境的服务器配置
// 逆向自 libMiniBaseEngine.dll 服务器域名表
var Servers = map[int]ServerSet{
	EnvDomestic: {
		ShequHTTP:  "shequ.mini1.cn",
		ShequHTTPS: "shequ.mini1.cn",
		WAPI:       "wapi.mini1.cn",
		Download:   "mdownload.mini1.cn",
	},
	EnvDomesticDev: {
		ShequHTTP:  "shequ-pre.mini1.cn",
		ShequHTTPS: "indevelop.mini1.cn",
		WAPI:       "wapi.mini1.cn",
		Download:   "mdownload.mini1.cn",
	},
	EnvOversea: {
		ShequHTTP:  "shequ.miniworldgame.com",
		ShequHTTPS: "shequ.miniworldgame.com",
		WAPI:       "wapi.miniworldgame.com",
		Download:   "hwmdownload.mini1.cn",
	},
	EnvOverseaDev: {
		ShequHTTP:  "shequ.miniworldplus.com",
		ShequHTTPS: "shequ.miniworldplus.com",
		WAPI:       "wapi.miniworldplus.com",
		Download:   "hwmdownload.mini1.cn",
	},
}

// 社区API端口 (逆向自 libMiniBaseEngine.dll)
const (
	ShequHTTPPort  = 8080
	ShequHTTPSPort = 8081
)

// DevServerIP 开发/调试直连IP (逆向自 libMiniBaseEngine.dll)
const (
	DevServerIP   = "120.24.64.132"
	DevServerPort = 8080
)

// DNS缓存IP (逆向自 iworld.cfg DnsCache节点)
const (
	DNSOpenRoom = "129.211.56.128" // openroom.mini1.cn
	DNSOperate  = "139.199.5.123"  // operate.mini1.cn
	DNSMap0     = "118.89.30.179"  // map0.mini1.cn
	DNSFriend   = "123.207.243.220" // friend.mini1.cn
)

// DNS解析候选IP (逆向自 libMiniBaseEngine.dll DnsMgr)
var DNSResolverIPs = []string{
	"119.3.38.56",
	"116.205.254.245",
	"124.71.120.6",
}

// 海外专用域名 (逆向自 iworld.cfg)
const (
	HWOpenRoomHost = "hwopenroom.mini1.cn"
	HWShequHost    = "hwshequ.mini1.cn"
	HWMailHost     = "hwmail.mini1.cn"
	HWFriendHost   = "hwfriend.mini1.cn"
)

// 热更新/版本检测 (逆向自 libMiniBaseEngine.dll)
const (
	UpdateServer    = "update.mini1.cn:13002"
	UpdatePkgPath   = "/miniw/patch_server"
	EngineAssets    = "engine.mini1.cn"
	VersionJSON     = "/game/version.json"
	VersionWebHost  = "mnweb.mini1.cn"
	HWVersionHost   = "mnweb.miniworldgame.com"
)
