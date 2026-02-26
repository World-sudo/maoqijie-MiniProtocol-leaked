// Package register 注册/登录接口
// 逆向发现:
//   - 注册路径: h5.mini1.cn/register/ (HTTPS, H5 WebView, GeeTest V4验证码)
//   - 登录路径: mnweb.mini1.cn/account/TextPwdLogin (HTTPS)
//   - 域名登录: <server>/miniw/ldap/auth?time=%d&auth=%s (POST, 内部LDAP)
//   - 回调: SDKLoginCallBack(%d, [[%s]], [[%s]], %d)  (status, uid, token, type)
//   - 回调: NativeCalledLoginManager([[DomainLoginResult]],[[...]])
//   - 自定义头: MN-AUTH, MN-TOKEN, MN-PAYLOAD
//   - JWT载荷: Uin, env, auth, ts, apiid, cltversion, src, deviceid, its, iat
//   - 认证前缀: switchAccountByAuthInfo_reg###<JWT>
//   - 签名密钥: #_php_miniw_2016_# (md5(uin+time+secret))
package register

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

// Client 注册/登录客户端
type Client struct {
	httpClient *httpc.Client
}

// NewClient 创建客户端
func NewClient(c *httpc.Client) *Client {
	return &Client{httpClient: c}
}

// LoginResult 登录响应
// 从 pcapng 抓包还原: 服务端返回 JSON 含 JWT 令牌
type LoginResult struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Token   string `json:"token"`
	Uin     int64  `json:"uin"`
	ErrType string `json:"err_type,omitempty"`
}

// DomainLogin 域名登录 (POST /miniw/ldap/auth)
// 逆向自 libiworld.dll: domainLogin: url=%s, postData=%s
// 参数: server - 服务器地址, uin - 用户ID, postData - 登录数据
// 回调: NativeCalledLoginManager([[DomainLoginResult]],[[httpcode,httpmsg]])
func (c *Client) DomainLogin(server string, uin int64, postData string) (int, string, error) {
	ts := time.Now().Unix()
	authSig := auth.Sign(uin, ts)

	u := fmt.Sprintf("%s%s?time=%d&auth=%s",
		server, config.DomainLoginPath, ts, authSig)

	req, err := http.NewRequest("POST", u,
		strings.NewReader(postData))
	if err != nil {
		return 0, "", fmt.Errorf("构造登录请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, resp.Status, nil
}

// TextPwdLogin 用户名密码登录 (POST /account/TextPwdLogin)
// 逆向自 LJ#261: uin + pwd 字段, 密码MD5后传输
// 回调: SDKLoginCallBack(%d, [[uid]], [[token]], %d)
func (c *Client) TextPwdLogin(account, password, deviceID string) (*LoginResult, error) {
	pwdMD5 := auth.MD5Password(password)

	params := url.Values{}
	params.Set("uin", account)
	params.Set("pwd", pwdMD5)
	params.Set("device_id", deviceID)

	u := fmt.Sprintf("https://%s%s",
		config.RegisterWebHost, config.TextPwdLoginPath)

	req, err := http.NewRequest("POST", u,
		strings.NewReader(params.Encode()))
	if err != nil {
		return nil, fmt.Errorf("构造登录请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	p := httpc.DefaultPayload(deviceID)
	httpc.InjectHeaders(req, "", "", p)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result LoginResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}

// TextPwdLoginURL 构造用户名密码登录URL
func TextPwdLoginURL() string {
	return fmt.Sprintf("https://%s%s",
		config.RegisterWebHost, config.TextPwdLoginPath)
}

// RegisterURL 构造注册页面URL
func RegisterURL() string {
	return fmt.Sprintf("https://%s%s",
		config.RegisterH5Host, config.RegisterPath)
}

// AuthAPIURL 构造通用认证API URL
// 逆向自 libiworld.dll: ?act=<op>&auth=%s&time=%u&uin=%u&s2t=%s&country=%s&lang=1
func AuthAPIURL(host, act string, uin int64, s2t string) string {
	ts := time.Now().Unix()
	authSig := auth.Sign(uin, ts)
	params := url.Values{}
	params.Set("act", act)
	params.Set("auth", authSig)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("uin", strconv.FormatInt(uin, 10))
	if s2t != "" {
		params.Set("s2t", s2t)
	}
	params.Set("country", "CN")
	params.Set("lang", "1")
	return fmt.Sprintf("https://%s?%s", host, params.Encode())
}

// BuildUpdateToken 构造更新检查token (md5签名)
func BuildUpdateToken(uin int64) string {
	ts := time.Now().Unix()
	return auth.Sign(uin, ts)
}

// BuildUpdateCheckURL 构造热更新检查URL
func BuildUpdateCheckURL(uin int64) string {
	token := BuildUpdateToken(uin)
	params := url.Values{}
	params.Set("app_ver", strconv.Itoa(config.CltVersion))
	params.Set("channel", strconv.Itoa(config.APIID))
	params.Set("env", "0")
	params.Set("os_type", "1")
	params.Set("token", token)
	params.Set("uin", strconv.FormatInt(uin, 10))
	return fmt.Sprintf("https://mwu-api-pre.mini1.cn/app_update/check_app_ver?%s",
		params.Encode())
}
