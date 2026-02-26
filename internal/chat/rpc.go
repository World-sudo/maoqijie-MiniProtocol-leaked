package chat

import (
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

// RPC 聊天RPC客户端
type RPC struct {
	client *httpc.Client
	cred   *auth.Credential
}

// NewRPC 创建RPC客户端
func NewRPC(client *httpc.Client, cred *auth.Credential) *RPC {
	return &RPC{client: client, cred: cred}
}

// Call 发起RPC调用
func (r *RPC) Call(body url.Values) ([]byte, error) {
	ts := time.Now().Unix()
	authSig := r.cred.ChatAuthAt(ts)

	params := url.Values{}
	params.Set("uid", strconv.FormatInt(r.cred.Uin, 10))
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("auth", authSig)
	params.Set("loginauth", r.cred.AuthString())
	params.Set("s2t", strconv.FormatInt(ts, 10))

	u := fmt.Sprintf("http://%s:%d%s?%s",
		config.ChatAllocHost, config.ChatAllocPort, config.ChatRPCPath,
		params.Encode())

	resp, err := r.client.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("RPC请求失败: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取RPC响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("RPC返回状态码: %d body=%s", resp.StatusCode, string(data))
	}
	return data, nil
}
