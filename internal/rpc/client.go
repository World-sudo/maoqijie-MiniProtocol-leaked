// Package rpc RPC通信客户端
// 逆向自 LJ#137 rpc_do_http_post + LJ#261 http_rpc
// 加密流程: plaintext → XXTEA(key) → Base64 → gsub(safe_7)
// 解密流程: response → Base64 → XXTEA → gzip解压
package rpc

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/crypto"
	"miniprotocol/internal/httpc"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// XXTEA密钥候选 (逆向自 patch_game_script.pkg)
// 在字节码中被拆分和字符重排列存储，运行时动态重组
// 片段: NpcVg21KdRWQ5t0y / sLer9mPoH4O3aZv / EBIwiTJCY8FjSkUbu / XM6lfqxAG7Dnzh
const XXTEAKeyCandidate = "NpcVg21KdRWQ5t0y"

// Client RPC客户端
type Client struct {
	httpClient *httpc.Client
	baseURL    string
	encKey     []byte
}

// NewClient 创建RPC客户端
// baseURL 格式: https://shequ.mini1.cn:8080 或使用测试IP
func NewClient(c *httpc.Client, baseURL string) *Client {
	return &Client{
		httpClient: c,
		baseURL:    baseURL,
		encKey:     []byte(XXTEAKeyCandidate),
	}
}

// SetKey 设置XXTEA密钥 (DH协商后调用)
func (c *Client) SetKey(key []byte) {
	c.encKey = key
}

// Call 发起RPC调用
// cmd - RPC命令名
// uin - 用户ID
// body - 请求体 (明文JSON或form-data)
func (c *Client) Call(cmd string, uin int64, body string) ([]byte, error) {
	encrypted := c.encryptZip([]byte(body))

	u := fmt.Sprintf("%s/_proxy?cmd=%s&uin=%d", c.baseURL, url.QueryEscape(cmd), uin)

	req, err := http.NewRequest("POST", u, strings.NewReader(encrypted))
	if err != nil {
		return nil, fmt.Errorf("构造RPC请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("RPC请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取RPC响应失败: %w", err)
	}

	return c.decryptUnzip(respBody)
}

// BuildSignedRPCURL 构建带签名的RPC URL
// 逆向自 LJ#261: _login_sign + gFunc_getmd5
func BuildSignedRPCURL(baseURL, act string, uin int64) string {
	ts := time.Now().Unix()
	authSig := auth.Sign(uin, ts)
	sign := auth.URLSignMD5(strconv.FormatInt(uin, 10))

	params := url.Values{}
	params.Set("act", act)
	params.Set("uin", strconv.FormatInt(uin, 10))
	params.Set("auth", authSig)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("sign", sign)
	params.Set("env", "0")
	params.Set("lang", "1")
	return baseURL + "?" + params.Encode()
}

// encryptZip XXTEA加密 + Base64编码
// 逆向自 LJ#137: encrypt_zip → xxtea → b64 → gsub(safe_7)
func (c *Client) encryptZip(data []byte) string {
	encrypted := crypto.XXTEAEncrypt(data, c.encKey)
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	// safe_7: URL安全替换 (+ → -, / → _, = → 去除)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", "_")
	encoded = strings.TrimRight(encoded, "=")
	return encoded
}

// decryptUnzip Base64解码 + XXTEA解密 + gzip解压
// 逆向自 LJ#137: decrypt_unzip → b64 → xxtea → unzip
func (c *Client) decryptUnzip(data []byte) ([]byte, error) {
	s := string(data)
	// 还原safe_7替换
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")
	// 补齐Base64 padding
	if m := len(s) % 4; m != 0 {
		s += strings.Repeat("=", 4-m)
	}

	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("base64解码失败: %w", err)
	}

	decrypted := crypto.XXTEADecrypt(decoded, c.encKey)

	// 尝试gzip解压 (部分响应未压缩)
	if len(decrypted) > 2 && decrypted[0] == 0x1f && decrypted[1] == 0x8b {
		return gzipDecompress(decrypted)
	}
	return decrypted, nil
}

func gzipDecompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("gzip reader: %w", err)
	}
	defer r.Close()
	return io.ReadAll(r)
}
