package chat

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

// AllocResponse 聊天节点分配响应
type AllocResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data AllocData `json:"data"`
}

// AllocData 分配数据
type AllocData struct {
	Token string `json:"token"`
	Host  string `json:"host"`
}

// Allocator 聊天节点分配器
type Allocator struct {
	client *httpc.Client
	cred   *auth.Credential
}

// NewAllocator 创建分配器
func NewAllocator(client *httpc.Client, cred *auth.Credential) *Allocator {
	return &Allocator{client: client, cred: cred}
}

// Alloc 分配聊天节点，返回 IM JWT 和 host
func (a *Allocator) Alloc() (*AllocResponse, error) {
	ts := time.Now().Unix()
	authSig := a.cred.ChatAuthAt(ts)
	uinStr := strconv.FormatInt(a.cred.Uin, 10)
	sign := auth.URLSignMD5(uinStr)

	form := url.Values{}
	form.Set("uid", uinStr)
	form.Set("time", strconv.FormatInt(ts, 10))
	form.Set("auth", authSig)
	form.Set("cltversion", strconv.Itoa(config.CltVersion))
	form.Set("apiid", strconv.Itoa(config.APIID))
	form.Set("env", "0")
	form.Set("s2t", "0")
	form.Set("country", "CN")
	form.Set("lang", "1")
	form.Set("sign", sign)

	u := fmt.Sprintf("http://%s:%d%s",
		config.ChatAllocHost, config.ChatAllocPort, config.ChatAllocPath)

	resp, err := a.client.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("聊天分配请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取分配响应失败: %w", err)
	}

	var result AllocResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析分配响应失败: %w (body=%s)", err, string(body))
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("聊天分配失败: code=%d msg=%s", result.Code, result.Msg)
	}

	// 保存 IM JWT 到凭证
	a.cred.SetIMJWT(result.Data.Token)
	return &result, nil
}
