// Package friend 好友系统
// 逆向推断自 iworld.cfg DnsCache: friend.mini1.cn (123.207.243.220)
// + shequ.mini1.cn 社区API的 act= 接口模式
// + libiworld.dll 好友管理相关字符串: FriendList, FriendOnline, FriendInfo
//
// API端点推断:
//   shequ.mini1.cn:8080/miniw/friend/?act=<op>
//   friend.mini1.cn/friend/<op>
//
// 签名方式与其他 shequ 接口一致: auth + sign + time + uin + env
package friend

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

// FriendInfo 好友基本信息
// 逆向推断自 libiworld.dll: FriendData 结构体字段
type FriendInfo struct {
	Uin       int64  `json:"uin"`
	NickName  string `json:"nick_name"`
	Level     int    `json:"level"`
	SkinID    int    `json:"skin_id"`
	IsOnline  bool   `json:"is_online"`
	LastLogin int64  `json:"last_login_time"`
	VIP       int    `json:"vip"`
}

// ListResponse 好友列表响应
type ListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total   int          `json:"total"`
		Friends []FriendInfo `json:"friends"`
	} `json:"data"`
}

// OnlineResponse 在线好友响应
type OnlineResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Online []int64 `json:"online_uins"`
	} `json:"data"`
}

// Service 好友服务
type Service struct {
	client *httpc.Client
	cred   *auth.Credential
}

// NewService 创建好友服务
func NewService(client *httpc.Client, cred *auth.Credential) *Service {
	return &Service{client: client, cred: cred}
}

// List 获取好友列表
// GET shequ.mini1.cn:8080/miniw/friend/?act=friend_list
func (s *Service) List() (*ListResponse, error) {
	params := s.baseParams()
	params.Set("act", "friend_list")

	u := fmt.Sprintf("http://%s:%d/miniw/friend/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[ListResponse](s.client, u)
}

// ListByPage 分页获取好友列表
// GET shequ.mini1.cn:8080/miniw/friend/?act=friend_list&page=<n>&size=<s>
func (s *Service) ListByPage(page, size int) (*ListResponse, error) {
	params := s.baseParams()
	params.Set("act", "friend_list")
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	u := fmt.Sprintf("http://%s:%d/miniw/friend/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[ListResponse](s.client, u)
}

// Online 查询在线好友列表
// GET shequ.mini1.cn:8080/miniw/friend/?act=friend_online
func (s *Service) Online() (*OnlineResponse, error) {
	params := s.baseParams()
	params.Set("act", "friend_online")

	u := fmt.Sprintf("http://%s:%d/miniw/friend/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[OnlineResponse](s.client, u)
}

func (s *Service) baseParams() url.Values {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(s.cred.Uin, 10)
	params := url.Values{}
	params.Set("uin", uinStr)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("auth", auth.Sign(s.cred.Uin, ts))
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
