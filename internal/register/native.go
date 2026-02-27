package register

// 原生登录模块
// 逆向自 MicroMiniNew.exe (反汇编 0x0045E600-0x0045F300)
// API: POST /login/auth_security (JSON body)
// 服务器: wskacchm.mini1.cn:14130
//
// 请求体结构 (已通过服务器验证):
//   公共字段: source, time(数字), auth, target
//   登录(target=auth): passwd_auth: {uin, passwd, apiid, DeviceID, cltversion}
//   注册(target=reg):  reg: {passwd, apiid, DeviceID, cltversion}
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

// NativeAuthResponse 原生认证通用响应
type NativeAuthResponse struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	StringID int    `json:"stringid,omitempty"`
	Data     any    `json:"data,omitempty"`
}

// RegisterResponse 注册成功响应 (code=0)
// authinfo/baseinfo 使用 AES-256-CBC 加密 (逆向自 MicroMiniNew.exe 0x0045A180)
// 解密方法在 decrypt.go 中
type RegisterResponse struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	AuthInfo string `json:"authinfo"`
	BaseInfo string `json:"baseinfo"`
	IV       string `json:"iv"`
}

// LoginResponse 登录成功响应 (code=0)
// 逆向自 MicroMiniNew.exe 0x0045D700 响应解析代码
// 解密方法在 decrypt.go 中
type LoginResponse struct {
	Code     int        `json:"code"`
	Msg      string     `json:"msg"`
	AuthInfo string     `json:"authinfo"`
	BaseInfo string     `json:"baseinfo"`
	IV       string     `json:"iv"`
	Data     *LoginData `json:"data,omitempty"`
}

// LoginData 登录成功时的data字段
type LoginData struct {
	Uin               int64  `json:"Uin"`
	Token             string `json:"token"`
	Sign              string `json:"sign"`
	IsLoginSafeVerify int    `json:"isloginsafeverify"`
}

// Login 原生密码登录 (target=auth, passwd_auth 方式)
// 已通过服务器验证: 返回 code=7012 表示密码错误, code=0 表示成功
func (c *NativeClient) Login(uin int64, passwd, deviceID string) (*LoginResponse, error) {
	ts := time.Now().Unix()

	req := map[string]any{
		"source": "mini_micro",
		"target": "auth",
		"time":   ts,
		"auth":   auth.NativeAuthSign("auth", fmt.Sprintf("%d", ts)),
		"passwd_auth": map[string]any{
			"uin":        uin,
			"passwd":     passwd,
			"apiid":      config.APIID,
			"DeviceID":   deviceID,
			"cltversion": config.CltVersion,
		},
	}
	return c.postLogin(req)
}

// Register 原生注册 (target=reg)
// 已通过服务器验证: 返回 code=0 msg="注册成功"
func (c *NativeClient) Register(passwd, deviceID string) (*RegisterResponse, error) {
	ts := time.Now().Unix()

	req := map[string]any{
		"source": "mini_micro",
		"target": "reg",
		"time":   ts,
		"auth":   auth.NativeAuthSign("reg", fmt.Sprintf("%d", ts)),
		"reg": map[string]any{
			"passwd":     passwd,
			"apiid":      config.APIID,
			"DeviceID":   deviceID,
			"cltversion": config.CltVersion,
		},
	}
	return c.postRegister(req)
}

// AuthInfoLogin 使用 authinfo_auth 令牌登录
// 注意: 此方式目前测试返回"验证失败,请重新登录"，可能需要额外处理
func (c *NativeClient) AuthInfoLogin(token, deviceID string) (*LoginResponse, error) {
	ts := time.Now().Unix()

	req := map[string]any{
		"source": "mini_micro",
		"target": "auth",
		"time":   ts,
		"auth":   auth.NativeAuthSign("auth", fmt.Sprintf("%d", ts)),
		"authinfo_auth": map[string]any{
			"token":      token,
			"apiid":      config.APIID,
			"DeviceID":   deviceID,
			"cltversion": config.CltVersion,
		},
	}
	return c.postLogin(req)
}

func (c *NativeClient) postLogin(payload any) (*LoginResponse, error) {
	respBody, err := c.doPost(payload)
	if err != nil {
		return nil, err
	}

	var result LoginResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(respBody))
	}
	return &result, nil
}

func (c *NativeClient) postRegister(payload any) (*RegisterResponse, error) {
	respBody, err := c.doPost(payload)
	if err != nil {
		return nil, err
	}

	var result RegisterResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(respBody))
	}
	return &result, nil
}

func (c *NativeClient) doPost(payload any) ([]byte, error) {
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
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
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
