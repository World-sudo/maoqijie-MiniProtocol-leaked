package room

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
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

	u := fmt.Sprintf("http://%s:%d%s?cmd=server_config&uin=%d&auth=%s",
		config.RoomHost, config.RoomPort, config.RoomPath,
		s.cred.Uin, authSig)

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
