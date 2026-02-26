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
	"fmt"
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
// 这些参数在运行时由DHKey对象初始化
// 目前使用标准的2048-bit MODP Group (RFC 3526 Group 14) 作为占位
// 实际参数需要从运行时内存dump提取
func DefaultDHParams() *DHParams {
	g := big.NewInt(2)
	// RFC 3526 Group 14 (2048-bit)
	p, _ := new(big.Int).SetString(
		"FFFFFFFFFFFFFFFFC90FDAA22168C234C4C6628B80DC1CD1"+
			"29024E088A67CC74020BBEA63B139B22514A08798E3404DD"+
			"EF9519B3CD3A431B302B0A6DF25F14374FE1356D6D51C245"+
			"E485B576625E7EC6F44C42E9A637ED6B0BFF5CB6F406B7ED"+
			"EE386BFB5A899FA5AE9F24117C4B1FE649286651ECE45B3D"+
			"C2007CB8A163BF0598DA48361C55D39A69163FA8FD24CF5F"+
			"83655D23DCA3AD961C62F356208552BB9ED529077096966D"+
			"670C354E4ABC9804F1746C08CA18217C32905E462E36CE3B"+
			"E39E772C180E86039B2783A2EC07A28FB5C55DF06F4C52C9"+
			"DE2BCBF6955817183995497CEA956AE515D2261898FA0510"+
			"15728E5A8AACAA68FFFFFFFFFFFFFFFF", 16)
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

	// TODO: 解析服务端返回的公钥B
	// 响应格式待进一步逆向确认
	return "", fmt.Errorf("DH交换响应解析待实现 (status=%d)", resp.StatusCode)
}
