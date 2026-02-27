package room

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

// Service 房间服务
type Service struct {
	client *httpc.Client
	cred   *auth.Credential
}

// NewService 创建房间服务
func NewService(client *httpc.Client, cred *auth.Credential) *Service {
	return &Service{client: client, cred: cred}
}

// GetConfig 获取房间配置
func (s *Service) GetConfig() (*Response, error) {
	ts := time.Now().Unix()
	authSig := s.cred.ChatAuthAt(ts)
	uinStr := strconv.FormatInt(s.cred.Uin, 10)
	sign := auth.URLSignMD5(uinStr)

	params := url.Values{}
	params.Set("cmd", "server_config")
	params.Set("uin", uinStr)
	params.Set("auth", authSig)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("env", "0")
	params.Set("s2t", "0")
	params.Set("country", "CN")
	params.Set("lang", "1")
	params.Set("sign", sign)

	u := fmt.Sprintf("http://%s:%d%s?%s",
		config.RoomHost, config.RoomPort, config.RoomPath,
		params.Encode())

	resp, err := s.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("请求房间配置失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析房间配置失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
