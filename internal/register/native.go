package register

// 原生登录模块
// 逆向自 MicroMiniNew.exe
// API: /login/auth_security
// 服务器: wskacchm.mini1.cn 或 120.24.63.165:14000
// 三种登录方式: TextPasswordLogin, DigitalPasswordLogin, CreateAccount
// 四种认证方式: passwd_auth, authinfo_auth, question_auth, phone_login_education

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/http"
	"net/url"
	"strings"
)

// NativeClient 原生登录客户端 (MicroMiniNew.exe 协议)
type NativeClient struct {
	httpClient *httpc.Client
	server     string
}

// NewNativeClient 创建原生登录客户端
// server 格式: https://wskacchm.mini1.cn:14130 或使用默认
func NewNativeClient(c *httpc.Client, server string) *NativeClient {
	if server == "" {
		server = fmt.Sprintf("https://%s:%d",
			config.ChannelHost, config.ChannelPortPre)
	}
	return &NativeClient{httpClient: c, server: server}
}

// NativeAuthRequest 原生认证请求参数
// 逆向自 MicroMiniNew.exe 字符串: Uin, token, sign, baseinfo, authinfo...
type NativeAuthRequest struct {
	Uin        string `json:"Uin"`
	Token      string `json:"token"`
	Sign       string `json:"sign"`
	BaseInfo   string `json:"baseinfo"`
	AuthInfo   string `json:"authinfo"`
	Passwd     string `json:"passwd"`
	CltVersion string `json:"cltversion"`
	DeviceID   string `json:"DeviceID"`
	APIID      string `json:"apiid"`
	Target     string `json:"target"`
	Source     string `json:"source"`
	Time       string `json:"time"`
}

// NativeAuthResponse 原生认证响应
type NativeAuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// AuthSecurity 调用原生认证接口
// POST /login/auth_security
func (c *NativeClient) AuthSecurity(params url.Values) (*NativeAuthResponse, error) {
	u := c.server + config.NativeAuthPath

	req, err := http.NewRequest("POST", u,
		strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("构造原生登录请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("原生登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result NativeAuthResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}

// SMSSend 发送短信验证码
// POST /sms/smssend/?phone=<phone>&checktype=2
func (c *NativeClient) SMSSend(phone string) (*NativeAuthResponse, error) {
	u := fmt.Sprintf("%s%s?phone=%s&checktype=%s",
		c.server, config.SMSSendPath, url.QueryEscape(phone), config.SMSCheckType)
	return c.doGet(u)
}

// SMSVerify 验证短信验证码
// GET /sms/smsverify/?phone=<phone>&id=461053&uin=<uin>
func (c *NativeClient) SMSVerify(phone, uin string) (*NativeAuthResponse, error) {
	u := fmt.Sprintf("%s%s?phone=%s&id=%s&uin=%s",
		c.server, config.SMSVerifyPath,
		url.QueryEscape(phone), config.SMSID, url.QueryEscape(uin))
	return c.doGet(u)
}

// EmailSend 发送邮箱验证码
// GET /email/emailsend/?email=<email>
func (c *NativeClient) EmailSend(email string) (*NativeAuthResponse, error) {
	u := fmt.Sprintf("%s%s?email=%s",
		c.server, config.EmailSendPath, url.QueryEscape(email))
	return c.doGet(u)
}

// EmailVerify 验证邮箱验证码
// GET /email/emailverify/?email=<email>&uin=<uin>&verify_code=<code>
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
