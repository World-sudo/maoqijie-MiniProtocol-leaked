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
	"strings"
	"time"
)

// ActionResponse 好友操作通用响应
type ActionResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// SearchResponse 搜索用户响应
type SearchResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Users []FriendInfo `json:"users"`
	} `json:"data"`
}

// RequestInfo 好友请求信息
type RequestInfo struct {
	Uin      int64  `json:"uin"`
	NickName string `json:"nick_name"`
	Message  string `json:"message"`
	Time     int64  `json:"time"`
}

// RequestListResponse 好友请求列表响应
type RequestListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Requests []RequestInfo `json:"requests"`
	} `json:"data"`
}

// Add 发送好友申请
// POST shequ.mini1.cn:8080/miniw/friend/?act=friend_add
func (s *Service) Add(targetUin int64, message string) (*ActionResponse, error) {
	form := s.baseForm()
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	if message != "" {
		form.Set("message", message)
	}
	return s.doPost("friend_add", form)
}

// Delete 删除好友
// POST shequ.mini1.cn:8080/miniw/friend/?act=friend_del
func (s *Service) Delete(targetUin int64) (*ActionResponse, error) {
	form := s.baseForm()
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	return s.doPost("friend_del", form)
}

// Accept 接受好友申请
// POST shequ.mini1.cn:8080/miniw/friend/?act=friend_accept
func (s *Service) Accept(targetUin int64) (*ActionResponse, error) {
	form := s.baseForm()
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	return s.doPost("friend_accept", form)
}

// Reject 拒绝好友申请
// POST shequ.mini1.cn:8080/miniw/friend/?act=friend_reject
func (s *Service) Reject(targetUin int64) (*ActionResponse, error) {
	form := s.baseForm()
	form.Set("target_uin", strconv.FormatInt(targetUin, 10))
	return s.doPost("friend_reject", form)
}

// Search 搜索用户 (按昵称或迷你号)
// GET shequ.mini1.cn:8080/miniw/friend/?act=friend_search&keyword=<kw>
func (s *Service) Search(keyword string) (*SearchResponse, error) {
	params := s.baseParams()
	params.Set("act", "friend_search")
	params.Set("keyword", keyword)

	u := fmt.Sprintf("http://%s:%d/miniw/friend/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[SearchResponse](s.client, u)
}

// Requests 获取好友申请列表
// GET shequ.mini1.cn:8080/miniw/friend/?act=friend_request_list
func (s *Service) Requests() (*RequestListResponse, error) {
	params := s.baseParams()
	params.Set("act", "friend_request_list")

	u := fmt.Sprintf("http://%s:%d/miniw/friend/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[RequestListResponse](s.client, u)
}

func (s *Service) baseForm() url.Values {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(s.cred.Uin, 10)
	form := url.Values{}
	form.Set("uin", uinStr)
	form.Set("time", strconv.FormatInt(ts, 10))
	form.Set("auth", auth.Sign(s.cred.Uin, ts))
	form.Set("sign", auth.URLSignMD5(uinStr))
	form.Set("env", "0")
	return form
}

func (s *Service) doPost(act string, form url.Values) (*ActionResponse, error) {
	params := url.Values{}
	params.Set("act", act)

	u := fmt.Sprintf("http://%s:%d/miniw/friend/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	resp, err := s.client.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("好友操作请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ActionResponse
	if err := parseJSON(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func parseJSON(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("解析响应失败: %w (body=%s)", err, string(data))
	}
	return nil
}
