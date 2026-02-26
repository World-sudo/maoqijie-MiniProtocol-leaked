package register

// SSO 登录模块
// 逆向自 sso.mini1.cn H5页面 (Vue.js SPA)
// 新版API: POST https://wapi.mini1.cn/login-service/api/web/v2/login
// 密码: Base64编码 (不是MD5!)
// 验证码: GeeTest V4, captcha_id = 57157b87c9788ae72be45a2c79c6dd1c
// 响应: code=0成功, code=7211需要验证码

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/http"
	"net/url"
	"strings"
)

// SSOClient SSO登录客户端
type SSOClient struct {
	httpClient *httpc.Client
}

// NewSSOClient 创建SSO客户端
func NewSSOClient(c *httpc.Client) *SSOClient {
	return &SSOClient{httpClient: c}
}

// SSOLoginRequest SSO登录请求体
// 逆向自 sso.mini1.cn JS: axios.post(newLogin, {uin, pwd, device, data, ver_id})
type SSOLoginRequest struct {
	Uin    int64  `json:"uin"`
	Pwd    string `json:"pwd"`
	Device string `json:"device"`
	Data   string `json:"data"`
	VerID  string `json:"ver_id"`
}

// SSOLoginResponse SSO登录响应
// code=0: 成功, code=7211: 需要验证码(data为ver_id)
type SSOLoginResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// GeeTestData GeeTest V4 验证数据
// 逆向自 sso.mini1.cn JS: captchaInstance.getValidate()
type GeeTestData struct {
	Platform      string `json:"platform"`
	Version       string `json:"version"`
	CaptchaID     string `json:"captcha_id"`
	LotNumber     string `json:"lot_number"`
	CaptchaOutput string `json:"captcha_output"`
	PassToken     string `json:"pass_token"`
	GenTime       string `json:"gen_time"`
}

// GeeTestWrapper GeeTest数据外层包装
// 逆向自 sso.mini1.cn JS: {type: "geetest", dataJson: JSON.stringify(n)}
type GeeTestWrapper struct {
	Type     string `json:"type"`
	DataJSON string `json:"dataJson"`
}

// SSOCodeNeedCaptcha SSO返回需要验证码的code
const SSOCodeNeedCaptcha = 7211

// Login SSO登录 (第一步: 不带验证码)
// 密码用Base64编码, 非MD5
func (c *SSOClient) Login(uin int64, password, deviceInfo string) (*SSOLoginResponse, error) {
	pwdB64 := base64.StdEncoding.EncodeToString([]byte(password))
	device := url.QueryEscape(deviceInfo)

	req := &SSOLoginRequest{
		Uin:    uin,
		Pwd:    pwdB64,
		Device: device,
		Data:   "",
		VerID:  "",
	}
	return c.doLogin(req)
}

// LoginWithCaptcha SSO登录 (第二步: 带GeeTest验证码)
func (c *SSOClient) LoginWithCaptcha(uin int64, password, deviceInfo, verID string, gt *GeeTestData) (*SSOLoginResponse, error) {
	pwdB64 := base64.StdEncoding.EncodeToString([]byte(password))
	device := url.QueryEscape(deviceInfo)

	// 构造验证数据: btoa(JSON.stringify({type:"geetest", dataJson:...}))
	gtJSON, _ := json.Marshal(gt)
	wrapper := &GeeTestWrapper{
		Type:     "geetest",
		DataJSON: string(gtJSON),
	}
	wrapperJSON, _ := json.Marshal(wrapper)
	data := base64.StdEncoding.EncodeToString(wrapperJSON)

	req := &SSOLoginRequest{
		Uin:    uin,
		Pwd:    pwdB64,
		Device: device,
		Data:   data,
		VerID:  verID,
	}
	return c.doLogin(req)
}

func (c *SSOClient) doLogin(req *SSOLoginRequest) (*SSOLoginResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	u := fmt.Sprintf("https://%s%s", config.SSOLoginAPI, config.SSOLoginPath)

	httpReq, err := http.NewRequest("POST", u, strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json;charset=UTF-8")
	httpReq.Header.Set("Origin", "https://"+config.SSOHost)
	httpReq.Header.Set("Referer", "https://"+config.SSOHost+"/")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("SSO登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result SSOLoginResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(respBody))
	}
	return &result, nil
}

// NeedsCaptcha 检查是否需要验证码
func (r *SSOLoginResponse) NeedsCaptcha() bool {
	return r.Code == SSOCodeNeedCaptcha
}

// VerID 获取验证码ver_id (当code=7211时)
func (r *SSOLoginResponse) VerID() string {
	if s, ok := r.Data.(string); ok {
		return s
	}
	return ""
}

// DefaultDevice 生成默认设备信息字符串
// 格式: {OS}-{浏览器/版本}
func DefaultDevice() string {
	return "Windows10-MiniWorldPC/" + config.AppVersion
}

// SSOLoginURL 构造SSO登录API地址
func SSOLoginURL() string {
	return fmt.Sprintf("https://%s%s", config.SSOLoginAPI, config.SSOLoginPath)
}

// SSOPageURL 构造SSO页面地址
func SSOPageURL() string {
	return fmt.Sprintf("https://%s/#/", config.SSOHost)
}
