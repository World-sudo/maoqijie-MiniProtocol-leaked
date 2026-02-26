package telemetry

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/url"
	"strings"
	"time"
)

// Reporter 遥测上报器
type Reporter struct {
	client *httpc.Client
	cred   *auth.Credential
}

// NewReporter 创建遥测上报器
func NewReporter(client *httpc.Client, cred *auth.Credential) *Reporter {
	return &Reporter{client: client, cred: cred}
}

// Report 上报事件列表到 tj3.mini1.cn
func (r *Reporter) Report(events []Event) error {
	return r.reportTo(config.TelemetryHost, events)
}

// ReportAlt 上报到备用遥测服务器 tj.mini1.cn
func (r *Reporter) ReportAlt(events []Event) error {
	return r.reportTo(config.TelemetryAltHost, events)
}

func (r *Reporter) reportTo(host string, events []Event) error {
	data, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("序列化事件失败: %w", err)
	}

	form := url.Values{}
	form.Set("json3", string(data))
	body := strings.NewReader(form.Encode())

	u := fmt.Sprintf("http://%s%s", host, config.TelemetryPath)
	resp, err := r.client.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		return fmt.Errorf("上报失败: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != 200 {
		return fmt.Errorf("上报返回状态码: %d", resp.StatusCode)
	}
	return nil
}

// LoginCheckConnectEvent 构建 login_check_connect 事件
// 逆向自 pcapng: 首次登录连接检查 (apiid=999)
func (r *Reporter) LoginCheckConnectEvent() Event {
	return Event{
		Age:             0,
		APIID:           999,
		Birthday:        0,
		ChannelAge:      0,
		ChannelRealname: 0,
		CltAPIID:        config.APIID,
		CltVersion:      config.CltVersion,
		Country:         "CN",
		CountryAuthPI:   r.cred.AuthString(),
		DeviceID:        r.cred.DeviceID,
		Env:             0,
		ID:              EventLoginCheckConnect,
		IP:              "0.0.0.0",
		IsAdult:         0,
		IsRealnamePass:  0,
		Language:        0,
		Log:             LogLoginCheckConnect,
		TS:              time.Now().Unix(),
		Uin:             r.cred.Uin,
	}
}

// LoginCheckEvent 构建 login_check 事件 (实名认证后)
// 逆向自 pcapng: 实名认证完成后的二次检查 (apiid=110, ip=真实IP)
func (r *Reporter) LoginCheckEvent(realIP string) Event {
	return Event{
		Age:            0,
		APIID:          config.APIID,
		Birthday:       "",
		CountryAuthPI:  0,
		DeviceID:       r.cred.DeviceID,
		IP:             realIP,
		IsRealnameConf: 1,
		IsRealnamePass: 1,
		Log:            LogLoginCheck,
		TS:             time.Now().Unix(),
		Uin:            r.cred.Uin,
	}
}
