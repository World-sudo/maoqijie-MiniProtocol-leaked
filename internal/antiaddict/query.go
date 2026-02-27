package antiaddict

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"strconv"
)

// QueryResponse 防沉迷查询响应
type QueryResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Status   int `json:"status"`
		Duration int `json:"duration"`
		Remain   int `json:"remain"`
	} `json:"data"`
}

// AntiAddictHost 防沉迷服务地址
const AntiAddictHost = "111.230.139.237:802"

// Client 防沉迷查询客户端
type Client struct {
	httpClient *httpc.Client
}

// NewClient 创建防沉迷查询客户端
func NewClient(c *httpc.Client) *Client {
	return &Client{httpClient: c}
}

// Query 查询防沉迷状态
// GET 111.230.139.237:802/antiaddiction.php?apiid=110&cmd_type=1&task_id=<uin>
func (c *Client) Query(uin int64) (*QueryResponse, error) {
	u := fmt.Sprintf("http://%s/antiaddiction.php?apiid=%d&cmd_type=1&task_id=%s",
		AntiAddictHost, config.APIID,
		strconv.FormatInt(uin, 10))

	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("防沉迷查询失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取防沉迷响应失败: %w", err)
	}

	var result QueryResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析防沉迷响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
