package mapapi

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

// ActionResponse 地图操作响应
type ActionResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Client 地图管理客户端
type Client struct {
	httpClient *httpc.Client
	cred       *auth.Credential
}

// NewClient 创建地图管理客户端
func NewClient(c *httpc.Client, cred *auth.Credential) *Client {
	return &Client{httpClient: c, cred: cred}
}

// RemoveMap 删除地图
// POST shequ.mini1.cn/miniw/map/?act=map_rm
func (c *Client) RemoveMap(mapID string) (*ActionResponse, error) {
	form := url.Values{}
	form.Set("map_id", mapID)
	return c.doAction("map_rm", form)
}

// UploadPreTime 预上传 (时间)
// POST shequ.mini1.cn/miniw/map/?act=upload_pre_time
func (c *Client) UploadPreTime(form url.Values) (*ActionResponse, error) {
	return c.doAction("upload_pre_time", form)
}

// UploadPrePlugin 预上传 (插件)
// POST shequ.mini1.cn/miniw/map/?act=upload_pre_plugin
func (c *Client) UploadPrePlugin(form url.Values) (*ActionResponse, error) {
	return c.doAction("upload_pre_plugin", form)
}

func (c *Client) doAction(act string, form url.Values) (*ActionResponse, error) {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	form.Set("uin", uinStr)
	form.Set("time", strconv.FormatInt(ts, 10))
	form.Set("auth", auth.Sign(c.cred.Uin, ts))
	form.Set("sign", auth.URLSignMD5(uinStr))
	form.Set("env", "0")

	params := url.Values{}
	params.Set("act", act)

	u := fmt.Sprintf("http://%s:%d/miniw/map/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	resp, err := c.httpClient.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("地图操作失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取地图响应失败: %w", err)
	}

	var result ActionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析地图响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
