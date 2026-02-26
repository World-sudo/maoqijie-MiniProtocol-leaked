package auth

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

// 签名密钥（从抓包推断，需实际验证）
const signSecret = "miniworld_chat_secret"

// Sign 生成MD5 auth签名: md5(uin + time + secret)
func Sign(uin int64, ts int64) string {
	raw := strconv.FormatInt(uin, 10) +
		strconv.FormatInt(ts, 10) +
		signSecret
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", hash)
}

// SignWithSecret 使用自定义secret生成签名
func SignWithSecret(uin int64, ts int64, secret string) string {
	raw := strconv.FormatInt(uin, 10) +
		strconv.FormatInt(ts, 10) +
		secret
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", hash)
}
