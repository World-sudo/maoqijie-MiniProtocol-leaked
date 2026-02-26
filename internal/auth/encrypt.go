package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"miniprotocol/internal/config"
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

// URLSign 生成URL签名: md5(uin+secret+ts)
// 逆向自 LJ#96: GetReqUrl / p_getmd5
func URLSign(uin int64, ts int64) string {
	return Sign(uin, ts)
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
