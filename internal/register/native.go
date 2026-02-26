package register

// 原生登录模块
// 逆向自 MicroMiniNew.exe (反汇编 0x0045E600-0x0045F300)
// API: POST /login/auth_security (JSON body, 不是 form-encoded)
// 服务器: wskacchm.mini1.cn:14130
//
// 请求体结构 (反汇编确认):
//   公共字段: source, time, auth, target, apiid(数字), DeviceID, cltversion(数字)
//   登录(target=auth): passwd_auth: {uin, passwd}
//   注册(target=reg):  reg: {passwd}   (注意: 注册服务当前不可用)
//
// 签名: auth = md5("source=mini_micro&target=<t>&time=<ts>" + serverSalt)
// 服务器使用的salt: 2ddb7619717147439c83ab022e9d4d38

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// NativeClient 原生登录客户端 (MicroMiniNew.exe 协议)
type NativeClient struct {
	httpClient *httpc.Client
	server     string
}

// NewNativeClient 创建原生登录客户端
// server 格式: https://wskacchm.mini1.cn:14130 或留空使用默认
func NewNativeClient(c *httpc.Client, server string) *NativeClient {
	if server == "" {
		server = fmt.Sprintf("https://%s:%d",
			config.ChannelHost, config.ChannelPortPre)
	}
	return &NativeClient{httpClient: c, server: server}
}

// PasswdAuthData 密码认证数据 (嵌套在 passwd_auth 字段中)
type PasswdAuthData struct {
	Uin    int64  `json:"uin"`
	Passwd string `json:"passwd"`
}

// RegData 注册数据 (嵌套在 reg 字段中)
type RegData struct {
	Passwd string `json:"passwd"`
}

// NativeAuthResponse 原生认证响应
type NativeAuthResponse struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	StringID int    `json:"stringid,omitempty"`
	Data     any    `json:"data,omitempty"`
}

// LoginData 登录成功时的响应数据
// 逆向自 MicroMiniNew.exe 0x0045D700 响应解析代码
type LoginData struct {
	AuthInfo          string `json:"authinfo"`
	BaseInfo          string `json:"baseinfo"`
	Uin               int64  `json:"Uin"`
	Token             string `json:"token"`
	Sign              string `json:"sign"`
	IsLoginSafeVerify int    `json:"isloginsafeverify"`
}

// Login 原生密码登录 (target=auth, passwd_auth 方式)
// 已通过服务器验证: apiid=110 有效, 返回 code=7012 表示密码错误
func (c *NativeClient) Login(uin int64, passwd, deviceID string) (*NativeAuthResponse, error) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	req := map[string]any{
		"source":     "mini_micro",
		"target":     "auth",
		"time":       ts,
		"auth":       auth.NativeAuthSign("auth", ts),
		"apiid":      config.APIID,
		"DeviceID":   deviceID,
		"cltversion": config.CltVersion,
		"passwd_auth": &PasswdAuthData{
			Uin:    uin,
			Passwd: passwd,
		},
	}
	return c.postJSON(req)
}

// Register 原生注册 (target=reg)
// 警告: 注册服务当前不可用，服务器对所有 apiid 值均返回 "apiid 不正确"
func (c *NativeClient) Register(passwd, deviceID string) (*NativeAuthResponse, error) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	req := map[string]any{
		"source":     "mini_micro",
		"target":     "reg",
		"time":       ts,
		"auth":       auth.NativeAuthSign("reg", ts),
		"apiid":      config.APIID,
		"DeviceID":   deviceID,
		"cltversion": config.CltVersion,
		"reg":        &RegData{Passwd: passwd},
	}
	return c.postJSON(req)
}

func (c *NativeClient) postJSON(payload any) (*NativeAuthResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	u := c.server + config.NativeAuthPath

	req, err := http.NewRequest("POST", u, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("原生登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result NativeAuthResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(respBody))
	}
	return &result, nil
}

// SMSSend 发送短信验证码
func (c *NativeClient) SMSSend(phone string) (*NativeAuthResponse, error) {
	u := fmt.Sprintf("%s%s?phone=%s&checktype=%s",
		c.server, config.SMSSendPath, url.QueryEscape(phone), config.SMSCheckType)
	return c.doGet(u)
}

// SMSVerify 验证短信验证码
func (c *NativeClient) SMSVerify(phone, uin string) (*NativeAuthResponse, error) {
	u := fmt.Sprintf("%s%s?phone=%s&id=%s&uin=%s",
		c.server, config.SMSVerifyPath,
		url.QueryEscape(phone), config.SMSID, url.QueryEscape(uin))
	return c.doGet(u)
}

// EmailSend 发送邮箱验证码
func (c *NativeClient) EmailSend(email string) (*NativeAuthResponse, error) {
	u := fmt.Sprintf("%s%s?email=%s",
		c.server, config.EmailSendPath, url.QueryEscape(email))
	return c.doGet(u)
}

// EmailVerify 验证邮箱验证码
func (c *NativeClient) EmailVerify(email, uin, code string) (*NativeAuthResponse, error) {
	u := fmt.Sprintf("%s%s?email=%s&uin=%s&verify_code=%s",
		c.server, config.EmailVerifyPath,
		url.QueryEscape(email), url.QueryEscape(uin), url.QueryEscape(code))
	return c.doGet(u)
}

func (c *NativeClient) doGet(u string) (*NativeAuthResponse, error) {
	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取失败: %w", err)
	}

	var result NativeAuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}

// NativeAuthURL 构造原生认证接口URL
func NativeAuthURL(server string) string {
	if server == "" {
		server = fmt.Sprintf("https://%s:%d",
			config.ChannelHost, config.ChannelPortPre)
	}
	return server + config.NativeAuthPath
}
