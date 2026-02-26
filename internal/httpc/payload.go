package httpc

import (
	"encoding/json"
	"miniprotocol/internal/config"
	"net/http"
	"strconv"
	"time"
)

// Payload MN-PAYLOAD 自定义头载荷
// 逆向自 libiworld.dll: MiniMutualParams (19个字段)
// 字段列表: version device_plat device_apn device_id device_ts
//
//	country_os country_ip nick uin_ts utc_offset zone
//	bind_mark client_type auth_ip auth_ts auth_mode
//	device_ip session_id apply_id
type Payload struct {
	Version    string `json:"version"`
	DevicePlat string `json:"device_plat"`
	DeviceAPN  string `json:"device_apn,omitempty"`
	DeviceID   string `json:"device_id"`
	DeviceTS   string `json:"device_ts"`
	CountryOS  string `json:"country_os"`
	CountryIP  string `json:"country_ip"`
	Nick       string `json:"nick,omitempty"`
	UinTS      string `json:"uin_ts,omitempty"`
	UTCOffset  string `json:"utc_offset"`
	Zone       string `json:"zone,omitempty"`
	BindMark   string `json:"bind_mark,omitempty"`
	ClientType string `json:"client_type"`
	AuthIP     string `json:"auth_ip,omitempty"`
	AuthTS     string `json:"auth_ts,omitempty"`
	AuthMode   string `json:"auth_mode,omitempty"`
	DeviceIP   string `json:"device_ip,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	ApplyID    string `json:"apply_id,omitempty"`
}

// SignParams 签名参数 (附加到MN-PAYLOAD或URL中)
// 逆向自 libiworld.dll: MiniMutualParams sign 区域
type SignParams struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	RequestID string `json:"requestId"`
	SignType  string `json:"sign_type"`
	SignVer   string `json:"sign_ver,omitempty"`
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
