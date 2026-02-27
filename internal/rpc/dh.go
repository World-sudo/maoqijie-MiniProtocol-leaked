package rpc

// DH密钥交换协议
// 逆向自 patch_game_script.pkg: DHKey + AesHelper + KeyAndIv
//
// 流程:
//   1. 客户端生成私钥 a, 计算公钥 A = g^a mod p
//   2. 通过 _proxy?act=L 发送 r_a=A 给服务端
//   3. 服务端返回公钥 B
//   4. 计算共享密钥 shared = B^a mod p
//   5. md5(shared) → 前16字节为AES Key, 后16字节为IV
//   6. 后续RPC通信使用此Key/IV进行XXTEA加密

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"miniprotocol/internal/httpc"
	"net/url"
)

// DHParams DH密钥交换参数
// 逆向自 patch_game_script.pkg @6845612: dh_g, hex_num
type DHParams struct {
	G *big.Int // 生成元
	P *big.Int // 大素数
}

// DefaultDHParams 默认DH参数
// 逆向提取的实际hex_num值 (非标准RFC参数)
func DefaultDHParams() *DHParams {
	g := big.NewInt(2)
	p, _ := new(big.Int).SetString(
		"19948998fdd5c7b3ade8e87b4a12e292b0b881a4e1ae97", 16)
	return &DHParams{G: g, P: p}
}

// DHExchange DH密钥交换器
type DHExchange struct {
	params     *DHParams
	privateKey *big.Int
	publicKey  *big.Int
}

// NewDHExchange 创建DH交换器
func NewDHExchange(params *DHParams) (*DHExchange, error) {
	// 生成256-bit随机私钥
	privBytes := make([]byte, 32)
	if _, err := rand.Read(privBytes); err != nil {
		return nil, fmt.Errorf("生成DH私钥失败: %w", err)
	}
	priv := new(big.Int).SetBytes(privBytes)

	// A = g^a mod p
	pub := new(big.Int).Exp(params.G, priv, params.P)

	return &DHExchange{
		params:     params,
		privateKey: priv,
		publicKey:  pub,
	}, nil
}

// PublicKey 获取公钥 (16进制字符串)
func (d *DHExchange) PublicKey() string {
	return d.publicKey.Text(16)
}

// DeriveKey 从对方公钥派生共享密钥 → Key + IV
// 逆向自 patch_game_script.pkg: KeyAndIv + md5t
func (d *DHExchange) DeriveKey(remotePublicHex string) (key, iv []byte, err error) {
	remotePub, ok := new(big.Int).SetString(remotePublicHex, 16)
	if !ok {
		return nil, nil, fmt.Errorf("无效的远端公钥: %s", remotePublicHex)
	}

	// shared = B^a mod p
	shared := new(big.Int).Exp(remotePub, d.privateKey, d.params.P)

	// md5(shared) → 分割为 key(16B) + iv(16B)
	hash := md5.Sum(shared.Bytes())
	key = hash[:8]
	iv = hash[8:]

	// 扩展到16字节 (XXTEA需要128-bit key)
	fullKey := make([]byte, 16)
	copy(fullKey, hash[:])
	return fullKey, iv, nil
}

// dhExchangeResponse 服务端DH交换响应
type dhExchangeResponse struct {
	Code int    `json:"code"`
	RB   string `json:"r_b"` // 服务端公钥B (hex)
}

// RequestDHExchange 通过 _proxy?act=L 发起DH密钥交换
// 逆向自 patch_game_script.pkg @6846035: _proxy?act=L&r_a=%s
func RequestDHExchange(client *httpc.Client, baseURL string, uin int64, publicKeyHex string) (string, error) {
	params := url.Values{}
	params.Set("act", "L")
	params.Set("B", fmt.Sprintf("%d", uin))
	params.Set("r_a", publicKeyHex)

	u := baseURL + "/_proxy?" + params.Encode()

	resp, err := client.Get(u)
	if err != nil {
		return "", fmt.Errorf("DH交换请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取DH响应失败: %w", err)
	}

	var result dhExchangeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析DH响应失败: %w (body=%s)", err, string(body))
	}

	if result.Code != 0 {
		return "", fmt.Errorf("DH交换失败: code=%d", result.Code)
	}

	if result.RB == "" {
		return "", fmt.Errorf("DH响应缺少r_b字段 (body=%s)", string(body))
	}

	return result.RB, nil
}
