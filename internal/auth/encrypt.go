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

// MD5Password 密码MD5哈希
// 逆向自 LJ#148: enryptPw / shouldNotMd5
func MD5Password(pwd string) string {
	hash := md5.Sum([]byte(pwd))
	return fmt.Sprintf("%x", hash)
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
