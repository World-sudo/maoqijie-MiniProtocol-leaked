package moderation

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// CheckResult 内容审核结果
type CheckResult struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Pass   bool   `json:"pass"`
		Reason string `json:"reason"`
	} `json:"data"`
}

// Checker 内容审核器
type Checker struct {
	client *httpc.Client
	cred   *auth.Credential
}

// NewChecker 创建内容审核器
func NewChecker(client *httpc.Client, cred *auth.Credential) *Checker {
	return &Checker{client: client, cred: cred}
}

// CheckText 文本内容审核
// POST shequ.mini1.cn/miniw/wordwall?act=checktxt2
func (c *Checker) CheckText(text, function string) (*CheckResult, error) {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	authSig := auth.Sign(c.cred.Uin, ts)
	sign := auth.URLSignMD5(uinStr)

	params := url.Values{}
	params.Set("act", "checktxt2")

	u := fmt.Sprintf("http://%s:%d/miniw/wordwall?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	form := url.Values{}
	form.Set("key", text)
	form.Set("time", strconv.FormatInt(ts, 10))
	form.Set("env", "0")
	form.Set("token", c.cred.LoginJWT)
	form.Set("source", "pc")
	form.Set("function", function)
	form.Set("auth", authSig)
	form.Set("uin", uinStr)
	form.Set("sign", sign)

	resp, err := c.client.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("内容审核请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取审核响应失败: %w", err)
	}

	var result CheckResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析审核结果失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
