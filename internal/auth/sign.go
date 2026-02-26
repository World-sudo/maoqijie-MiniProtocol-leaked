package auth

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

// 签名密钥（从 libiworld.dll 逆向提取）
// 格式: md5("%s%d%s" % (uin, time, secret))
const signSecret = "#_php_miniw_2016_#"

// URLSignSalt 通用API签名salt (逆向自 LJ#096 p_getmd5)
// 用途: 所有API请求的URL &sign= 参数
// 签名方式: md5("uin=" + uin + salt)
const URLSignSalt = "7vbrtqudwf#z6pb&c6m4%j#zujz7g72q#mbW5Wh7@CILChaW^6RqvRJtkntsie3"

// LoginSignSalt 登录专用签名salt (逆向自 LJ#261 login_sign)
// 用途: 登录RPC的POST body签名
// 签名方式: md5(拼接参数 + salt)
const LoginSignSalt = "cb86b2a814cd477703073fb440386562"

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

// URLSignMD5 生成通用API URL签名 (LJ#096 p_getmd5)
// 用于附加到请求URL的 &sign= 参数
func URLSignMD5(uin string) string {
	raw := "uin=" + uin + URLSignSalt
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", hash)
}

// LoginSign 生成登录RPC签名 (LJ#261 login_sign)
// 用于登录POST请求体的签名字段
func LoginSign(params string) string {
	raw := params + LoginSignSalt
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", hash)
}

// NativeAuthSign 生成原生认证接口的auth字段
// 格式: md5("source=mini_micro&target=<target>&time=<ts>" + NativeServerSalt)
// 逆向自 MicroMiniNew.exe 0x0045E750 签名构造代码
// 注意: JSON字段名为 "auth"，不是 "sign"（已通过服务器验证确认）
func NativeAuthSign(target, timestamp string) string {
	raw := "source=mini_micro&target=" + target + "&time=" + timestamp +
		"2ddb7619717147439c83ab022e9d4d38"
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%x", hash)
}

// MD5Str 通用MD5字符串哈希
func MD5Str(s string) string {
	hash := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", hash)
}
