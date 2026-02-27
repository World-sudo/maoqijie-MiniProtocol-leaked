package credit

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"strconv"
)

// ScoreResponse 信用分查询响应
type ScoreResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Score     int  `json:"score"`
		Limited   bool `json:"limited"`
		ForbidEnd int  `json:"forbid_end"`
	} `json:"data"`
}

// Client 信用分查询客户端
type Client struct {
	httpClient *httpc.Client
}

// NewClient 创建信用分查询客户端
func NewClient(c *httpc.Client) *Client {
	return &Client{httpClient: c}
}

// QueryScore 查询用户信用分
// GET credit-api.mini1.cn/api/v1/action_limit/user?user_uin=xxx
func (c *Client) QueryScore(uin int64) (*ScoreResponse, error) {
	u := fmt.Sprintf("https://%s/api/v1/action_limit/user?user_uin=%s",
		config.CreditAPI, strconv.FormatInt(uin, 10))

	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("信用分查询失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取信用分响应失败: %w", err)
	}

	var result ScoreResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析信用分失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
