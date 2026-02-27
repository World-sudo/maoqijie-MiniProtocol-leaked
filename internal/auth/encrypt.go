package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"miniprotocol/internal/config"
	"strings"
)

// EncryptBDInfo AES-128-CBC加密bdinfo参数
// 逆向自 LJ#137 m_security模块: aes_encrypt
func EncryptBDInfo(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(config.AESKey))
	if err != nil {
		return "", fmt.Errorf("aes cipher: %w", err)
	}

	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	ct := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, []byte(config.AESIV))
	mode.CryptBlocks(ct, padded)

	return base64.StdEncoding.EncodeToString(ct), nil
}

// MD5Password 密码MD5哈希 (旧版TextPwdLogin用)
// 逆向自 LJ#148: enryptPw / shouldNotMd5
func MD5Password(pwd string) string {
	hash := md5.Sum([]byte(pwd))
	return fmt.Sprintf("%x", hash)
}

// Base64Password 密码Base64编码 (新版SSO登录用)
// 逆向自 sso.mini1.cn JS: h.encode() = 标准Base64
func Base64Password(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(pwd))
}

// NativeBase64Alphabet MicroMiniNew.exe自定义Base64字母表
// 逆向自 MicroMiniNew.exe 二进制字符串
const NativeBase64Alphabet = "Vg21WQ5KdRt0yNpcr9m4O3PoHaZvsLeCY8FjSwiTkUbuEBIJlAG7fqXM6xDnzh-;"

// NativeBase64Encode 使用自定义字母表的Base64编码
func NativeBase64Encode(data []byte) string {
	enc := base64.NewEncoding(NativeBase64Alphabet)
	return enc.EncodeToString(data)
}

// NativeBase64Decode 使用自定义字母表的Base64解码
func NativeBase64Decode(s string) ([]byte, error) {
	enc := base64.NewEncoding(NativeBase64Alphabet)
	return enc.DecodeString(s)
}

// DecryptBDInfo AES-128-CBC解密bdinfo参数
func DecryptBDInfo(cipherB64 string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(cipherB64)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	block, err := aes.NewCipher([]byte(config.AESKey))
	if err != nil {
		return "", fmt.Errorf("aes cipher: %w", err)
	}
	mode := cipher.NewCBCDecrypter(block, []byte(config.AESIV))
	mode.CryptBlocks(ct, ct)
	return string(pkcs7Unpad(ct)), nil
}

// DecryptNativeResponse 解密原生登录响应的 authinfo/baseinfo 字段
// 逆向自 MicroMiniNew.exe 0x0045A180 解密函数
//
// 完整流程 (从 0x0045D450 响应解析函数):
//  1. 从响应 JSON 中取出 authinfo/baseinfo (加密字符串) 和 iv (数字字符串)
//  2. 从 iv 字段解析数字 N (0x0045A1DE: 从末尾反向解析十进制数)
//  3. 字符重排: rotated = encrypted[len-N:] + encrypted[:len-N] (循环右移N位)
//  4. 标准 Base64 解码 (字母表: A-Za-z0-9+/, 初始化于 0x0064C420)
//  5. AES-256-CBC 解密 (0x0044C2B0 -> 0x0044B7D0 密钥设置 + 0x0044BBF0 解密)
//  6. PKCS7 去除填充
//
// authinfo 解密后包含: {"Uin": <int>, "token": "<string>", "sign": "<string>"}
// baseinfo 解密后包含: {"LastLoginTime": <int>, "isloginsafeverify": <bool>}
func DecryptNativeResponse(encrypted, iv string) (string, error) {
	// 步骤1: 从 iv 字段解析数字 N
	n := parseIVNumber(iv)
	if n <= 0 || n > len(encrypted) {
		return "", fmt.Errorf("invalid iv number: %d (encrypted len=%d)", n, len(encrypted))
	}

	// 步骤2: 字符重排 (循环右移N位)
	// 逆向自 0x0045A246/0x0045A25D: substr(0,len-N) + substr(len-N,N) -> concat(后半,前半)
	splitPos := len(encrypted) - n
	rotated := encrypted[splitPos:] + encrypted[:splitPos]

	// 步骤3: 标准 Base64 解码
	cipherBytes, err := base64.StdEncoding.DecodeString(rotated)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	// 步骤4: AES-256-CBC 解密
	block, err := aes.NewCipher([]byte(config.NativeRespAESKey))
	if err != nil {
		return "", fmt.Errorf("aes cipher: %w", err)
	}
	if len(cipherBytes)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext not aligned: %d bytes", len(cipherBytes))
	}
	mode := cipher.NewCBCDecrypter(block, []byte(config.NativeRespAESIV))
	mode.CryptBlocks(cipherBytes, cipherBytes)

	// 步骤5: PKCS7 去除填充
	plaintext := pkcs7Unpad(cipherBytes)
	return string(plaintext), nil
}

// parseIVNumber 从 iv 字符串解析数字 N
// 逆向自 0x0045A1DE-0x0045A228:
// 从字符串末尾(右端)向前遍历,遇到数字则累加: N += (char-'0') * placeValue; placeValue *= 10
// 遇到非数字时跳过但 placeValue 不重置
func parseIVNumber(iv string) int {
	n := 0
	placeValue := 1
	for i := len(iv) - 1; i >= 0; i-- {
		c := iv[i]
		if c >= '0' && c <= '9' {
			n += int(c-'0') * placeValue
			placeValue *= 10
		}
	}
	return n
}

// URLSign 生成URL签名: md5(uin+secret+ts)
// 逆向自 LJ#96: GetReqUrl / p_getmd5
func URLSign(uin int64, ts int64) string {
	return Sign(uin, ts)
}

// EnctyptW EnctyptW加密 (注意: 原始代码拼写错误 Enctypt)
// 逆向自 lj_098: privateKey + EnctyptW + md5Str + URL编码
// 用途: 特殊参数加密, 使用 DomainLoginHash 作为 privateKey
func EnctyptW(input, privateKey string) string {
	raw := input + privateKey
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", hash)
}

// URLEncodeBytes 手动URL编码 (逆向自 lj_098: format("%%%02X", byte))
func URLEncodeBytes(data []byte) string {
	var b strings.Builder
	for _, c := range data {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') ||
			(c >= '0' && c <= '9') || c == '.' || c == '-' || c == '_' || c == ' ' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	for i := len(data); i < len(padded); i++ {
		padded[i] = byte(padding)
	}
	return padded
}

func pkcs7Unpad(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return data
	}
	return data[:len(data)-padding]
}
