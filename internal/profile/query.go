// Package profile 用户资料查询与编辑
// 逆向推断自 libiworld.dll: UserInfo/UserCard 相关字符串
// + shequ.mini1.cn 社区API的 /miniw/user/?act= 接口模式
// + baseinfo 解密结构中的 NickName/SkinID/Level 等字段
//
// API端点:
//   shequ.mini1.cn:8080/miniw/user/?act=get_user_info
//   shequ.mini1.cn:8080/miniw/user/?act=get_user_card
//   shequ.mini1.cn:8081 (HTTPS 端口)
package profile

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/url"
	"strconv"
	"time"
)

// UserInfo 用户完整资料
// 逆向推断自 libiworld.dll 字段 + baseinfo 解密结构
type UserInfo struct {
	Uin        int64  `json:"uin"`
	NickName   string `json:"nick_name"`
	Sign       string `json:"sign"`
	Level      int    `json:"level"`
	Exp        int    `json:"exp"`
	VIP        int    `json:"vip"`
	SkinID     int    `json:"skin_id"`
	Gender     int    `json:"gender"`
	Birthday   string `json:"birthday"`
	Country    string `json:"country"`
	MiniCoin   int    `json:"minicoin"`
	MiniBean   int    `json:"minibean"`
	CreateTime int64  `json:"create_time"`
	MapCount   int    `json:"map_count"`
	FansCount  int    `json:"fans_count"`
	FollowCnt  int    `json:"follow_count"`
	FriendCnt  int    `json:"friend_count"`
	IsOnline   bool   `json:"is_online"`
}

// InfoResponse 用户信息查询响应
type InfoResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data *UserInfo `json:"data"`
}

// UserCard 用户名片 (简略)
type UserCard struct {
	Uin      int64  `json:"uin"`
	NickName string `json:"nick_name"`
	Level    int    `json:"level"`
	VIP      int    `json:"vip"`
	SkinID   int    `json:"skin_id"`
	Sign     string `json:"sign"`
	IsOnline bool   `json:"is_online"`
}

// CardResponse 用户名片响应
type CardResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data *UserCard `json:"data"`
}

// Client 用户资料客户端
type Client struct {
	httpClient *httpc.Client
	cred       *auth.Credential
}

// NewClient 创建用户资料客户端
func NewClient(c *httpc.Client, cred *auth.Credential) *Client {
	return &Client{httpClient: c, cred: cred}
}

// GetInfo 查询用户完整资料
// GET shequ.mini1.cn:8080/miniw/user/?act=get_user_info&target_uin=<uin>
func (c *Client) GetInfo(targetUin int64) (*InfoResponse, error) {
	params := c.baseParams()
	params.Set("act", "get_user_info")
	params.Set("target_uin", strconv.FormatInt(targetUin, 10))

	u := fmt.Sprintf("http://%s:%d/miniw/user/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[InfoResponse](c.httpClient, u)
}

// GetMyInfo 查询自己的完整资料
func (c *Client) GetMyInfo() (*InfoResponse, error) {
	return c.GetInfo(c.cred.Uin)
}

// GetCard 查询用户名片 (简略信息)
// GET shequ.mini1.cn:8080/miniw/user/?act=get_user_card&target_uin=<uin>
func (c *Client) GetCard(targetUin int64) (*CardResponse, error) {
	params := c.baseParams()
	params.Set("act", "get_user_card")
	params.Set("target_uin", strconv.FormatInt(targetUin, 10))

	u := fmt.Sprintf("http://%s:%d/miniw/user/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[CardResponse](c.httpClient, u)
}

func (c *Client) baseParams() url.Values {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	params := url.Values{}
	params.Set("uin", uinStr)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("auth", auth.Sign(c.cred.Uin, ts))
	params.Set("sign", auth.URLSignMD5(uinStr))
	params.Set("env", "0")
	params.Set("country", "CN")
	params.Set("lang", "1")
	return params
}

func doGet[T any](client *httpc.Client, u string) (*T, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
