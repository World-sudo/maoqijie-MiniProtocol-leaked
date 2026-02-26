package httpc

import (
	"encoding/json"
	"miniprotocol/internal/config"
	"net/http"
	"strconv"
	"time"
)

// Payload MN-PAYLOAD 自定义头载荷
// 逆向自 libiworld.dll: MN-AUTH MN-TOKEN MN-PAYLOAD
// 字段: version device_plat device_id device_ts country_os
//       nick uin_ts auth_mode session_id 等
type Payload struct {
	Version    string `json:"version"`
	DevicePlat string `json:"device_plat"`
	DeviceID   string `json:"device_id"`
	DeviceTS   string `json:"device_ts"`
	CountryOS  string `json:"country_os"`
	CountryIP  string `json:"country_ip"`
	Nick       string `json:"nick,omitempty"`
	UinTS      string `json:"uin_ts,omitempty"`
	UTCOffset  string `json:"utc_offset"`
	ClientType string `json:"client_type"`
	AuthMode   string `json:"auth_mode,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
}

// DefaultPayload 构造默认的MN-PAYLOAD
func DefaultPayload(deviceID string) *Payload {
	_, offset := time.Now().Zone()
	return &Payload{
		Version:    strconv.Itoa(config.CltVersion),
		DevicePlat: "1",
		DeviceID:   deviceID,
		DeviceTS:   strconv.FormatInt(time.Now().Unix(), 10),
		CountryOS:  "CN",
		CountryIP:  "CN",
		UTCOffset:  strconv.Itoa(offset / 3600),
		ClientType: "pc",
	}
}

// InjectHeaders 注入 MN-AUTH, MN-TOKEN, MN-PAYLOAD 到请求头
func InjectHeaders(req *http.Request, authVal, tokenVal string, p *Payload) {
	if authVal != "" {
		req.Header.Set(config.HeaderAuth, authVal)
	}
	if tokenVal != "" {
		req.Header.Set(config.HeaderToken, tokenVal)
	}
	if p != nil {
		data, _ := json.Marshal(p)
		req.Header.Set(config.HeaderPayload, string(data))
	}
}
