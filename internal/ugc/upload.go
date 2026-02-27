package ugc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// UploadResponse 上传响应
type UploadResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		URL   string `json:"url"`
		ResID string `json:"res_id"`
	} `json:"data"`
}

// ResIDResponse getResID 响应
type ResIDResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		ResID string `json:"res_id"`
	} `json:"data"`
}

// ResActionResponse 资源操作响应
type ResActionResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// Client UGC资源管理客户端
type Client struct {
	httpClient *httpc.Client
	cred       *auth.Credential
}

// NewClient 创建UGC客户端
func NewClient(c *httpc.Client, cred *auth.Credential) *Client {
	return &Client{httpClient: c, cred: cred}
}

// Upload 上传资源文件
// POST shequ.mini1.cn/v1/upload (multipart)
func (c *Client) Upload(filename string, data []byte) (*UploadResponse, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("创建multipart失败: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("写入文件数据失败: %w", err)
	}

	c.writeSignFields(w)
	w.Close()

	u := fmt.Sprintf("http://%s:%d/v1/upload",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort)

	resp, err := c.httpClient.Post(u, w.FormDataContentType(), &buf)
	if err != nil {
		return nil, fmt.Errorf("上传请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取上传响应失败: %w", err)
	}

	var result UploadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析上传响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}

// GetResID 获取资源ID
// GET shequ.mini1.cn/v1/getResID
func (c *Client) GetResID() (*ResIDResponse, error) {
	u := c.buildSignedURL("/v1/getResID")

	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, fmt.Errorf("获取ResID失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取ResID响应失败: %w", err)
	}

	var result ResIDResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析ResID失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}

// ResAction 执行资源操作 (upload/add_res/update_res)
// POST shequ.mini1.cn/miniw/res/v3?act=<action>
func (c *Client) ResAction(act string, form url.Values) (*ResActionResponse, error) {
	c.injectSignParams(form)

	params := url.Values{}
	params.Set("act", act)

	u := fmt.Sprintf("http://%s:%d/miniw/res/v3?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	resp, err := c.httpClient.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("资源操作失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取资源操作响应失败: %w", err)
	}

	var result ResActionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析资源操作响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}

func (c *Client) writeSignFields(w *multipart.Writer) {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	w.WriteField("uin", uinStr)
	w.WriteField("time", strconv.FormatInt(ts, 10))
	w.WriteField("auth", auth.Sign(c.cred.Uin, ts))
	w.WriteField("sign", auth.URLSignMD5(uinStr))
	w.WriteField("env", "0")
}

func (c *Client) buildSignedURL(path string) string {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	params := url.Values{}
	params.Set("uin", uinStr)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("auth", auth.Sign(c.cred.Uin, ts))
	params.Set("sign", auth.URLSignMD5(uinStr))
	params.Set("env", "0")
	return fmt.Sprintf("http://%s:%d%s?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, path, params.Encode())
}

func (c *Client) injectSignParams(form url.Values) {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	form.Set("uin", uinStr)
	form.Set("time", strconv.FormatInt(ts, 10))
	form.Set("auth", auth.Sign(c.cred.Uin, ts))
	form.Set("sign", auth.URLSignMD5(uinStr))
	form.Set("env", "0")
}
