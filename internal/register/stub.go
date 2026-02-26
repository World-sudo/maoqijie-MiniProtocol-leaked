// Package register 注册/登录接口
// 逆向发现:
//   - 注册路径: mnweb.mini1.cn/register/ (HTTPS, 需GeeTest V4验证码)
//   - 登录路径: mnweb.mini1.cn/account/TextPwdLogin (HTTPS)
//   - 域名登录: <server>/miniw/ldap/auth?time=%d&auth=%s (POST)
//   - 自定义头: MN-AUTH, MN-TOKEN, MN-PAYLOAD
//   - Payload字段: version, device_plat, device_id, device_ts,
//     country_os, country_ip, nick, uin_ts, auth_mode, session_id 等
//   - 登录回调: NativeCalledLoginManager([[DomainLoginResult]])
//   - 签名密钥: #_php_miniw_2016_# (md5(uin+time+secret))
package register

import (
	"fmt"
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

// DomainLogin 域名登录 (POST /miniw/ldap/auth)
// 参数: server - 服务器地址, uin - 用户ID, postData - 登录数据
// 返回: 响应体, 错误
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
