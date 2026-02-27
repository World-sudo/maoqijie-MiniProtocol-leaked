// Package social 社交代理
// 逆向自 config 中的 SocialProxyCmd = "/social_proxy" + RPCProxyPath = "/_proxy"
// 通过 shequ.mini1.cn:8080/_proxy?cmd=/social_proxy 发起社交操作
//
// 社交操作类型推断自 libiworld.dll 字符串:
//   follow, unfollow, like, unlike, fans_list, follow_list
//
// 签名方式复用 auth.Sign + auth.URLSignMD5
package social

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

// ProxyResponse 社交代理通用响应
type ProxyResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// UserBrief 用户简要信息 (粉丝/关注列表项)
type UserBrief struct {
	Uin      int64  `json:"uin"`
	NickName string `json:"nick_name"`
	Level    int    `json:"level"`
	VIP      int    `json:"vip"`
	SkinID   int    `json:"skin_id"`
}

// ListResponse 列表响应
type ListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total int         `json:"total"`
		List  []UserBrief `json:"list"`
	} `json:"data"`
}

// Client 社交代理客户端
type Client struct {
	httpClient *httpc.Client
	cred       *auth.Credential
}

// NewClient 创建社交代理客户端
func NewClient(c *httpc.Client, cred *auth.Credential) *Client {
	return &Client{httpClient: c, cred: cred}
}

// Follow 关注用户
// POST _proxy?cmd=/social_proxy, body: act=follow&target_uin=<uin>
func (c *Client) Follow(targetUin int64) (*ProxyResponse, error) {
	form := c.baseForm()
	form.Set("act", "follow")
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	return c.doProxyPost(form)
}

// Unfollow 取消关注
func (c *Client) Unfollow(targetUin int64) (*ProxyResponse, error) {
	form := c.baseForm()
	form.Set("act", "unfollow")
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	return c.doProxyPost(form)
}

// Like 点赞用户/作品
// resType: user=用户, map=地图作品
func (c *Client) Like(targetUin int64, resType string, resID string) (*ProxyResponse, error) {
	form := c.baseForm()
	form.Set("act", "like")
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	form.Set("res_type", resType)
	if resID != "" {
		form.Set("res_id", resID)
	}
	return c.doProxyPost(form)
}

// Unlike 取消点赞
func (c *Client) Unlike(targetUin int64, resType string, resID string) (*ProxyResponse, error) {
	form := c.baseForm()
	form.Set("act", "unlike")
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	form.Set("res_type", resType)
	if resID != "" {
		form.Set("res_id", resID)
	}
	return c.doProxyPost(form)
}

// FansList 获取粉丝列表
// GET _proxy?cmd=/social_proxy&act=fans_list&target_uin=<uin>&page=<p>&size=<s>
func (c *Client) FansList(targetUin int64, page, size int) (*ListResponse, error) {
	params := c.baseParams()
	params.Set("act", "fans_list")
	params.Set("target_uin", strconv.FormatInt(targetUin, 10))
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))
	return c.doProxyGet(params)
}

// FollowList 获取关注列表
func (c *Client) FollowList(targetUin int64, page, size int) (*ListResponse, error) {
	params := c.baseParams()
	params.Set("act", "follow_list")
	params.Set("target_uin", strconv.FormatInt(targetUin, 10))
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))
	return c.doProxyGet(params)
}

func (c *Client) baseForm() url.Values {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	form := url.Values{}
	form.Set("uin", uinStr)
	form.Set("time", strconv.FormatInt(ts, 10))
	form.Set("auth", auth.Sign(c.cred.Uin, ts))
	form.Set("sign", auth.URLSignMD5(uinStr))
	form.Set("env", "0")
	return form
}

func (c *Client) baseParams() url.Values {
	return c.baseForm()
}

func (c *Client) proxyURL() string {
	return fmt.Sprintf("http://%s:%d%s?cmd=%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort,
		config.RPCProxyPath,
		url.QueryEscape(config.SocialProxyCmd))
}

func (c *Client) doProxyPost(form url.Values) (*ProxyResponse, error) {
	resp, err := c.httpClient.Post(c.proxyURL(),
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("社交代理请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ProxyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}

func (c *Client) doProxyGet(params url.Values) (*ListResponse, error) {
	u := c.proxyURL() + "&" + params.Encode()

	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("社交代理请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
