package httpc

import (
	"fmt"
	"io"
	"miniprotocol/internal/config"
	"net/http"
	"time"
)

// Client 公共HTTP客户端，自动注入UA
type Client struct {
	inner *http.Client
}

// New 创建HTTP客户端
func New() *Client {
	return &Client{
		inner: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Do 发送请求，自动注入User-Agent
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", config.UserAgent)
	}
	return c.inner.Do(req)
}

// Get 便捷GET请求
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("构造GET请求失败: %w", err)
	}
	return c.Do(req)
}

// Post 便捷POST请求
func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("构造POST请求失败: %w", err)
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}
