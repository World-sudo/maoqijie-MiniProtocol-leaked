package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// LoginClaims 登录JWT载荷
type LoginClaims struct {
	Uin        int64  `json:"Uin"`
	Env        int    `json:"env"`
	Auth       string `json:"auth"`
	TS         int64  `json:"ts"`
	APIID      int    `json:"apiid"`
	CltVersion int    `json:"cltversion"`
	Src        string `json:"src"`
	DeviceID   string `json:"deviceid"`
	ITS        int64  `json:"its"`
	IAT        int64  `json:"iat"`
}

// IMClaims IM聊天JWT载荷
type IMClaims struct {
	Uin  string `json:"uin"`
	Time int64  `json:"time"`
	Flag int    `json:"flag"`
	Exp  int64  `json:"exp"`
	Iss  string `json:"iss"`
}

// jwtHeader HS256 JWT头
var jwtHeader = map[string]string{
	"alg": "HS256",
	"typ": "JWT",
}

// BuildLoginJWT 构造登录JWT
func BuildLoginJWT(claims *LoginClaims, secret string) (string, error) {
	headerJSON, _ := json.Marshal(jwtHeader)
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("序列化claims失败: %w", err)
	}

	headerB64 := base64URLEncode(headerJSON)
	claimsB64 := base64URLEncode(claimsJSON)
	sigInput := headerB64 + "." + claimsB64

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(sigInput))
	sig := base64URLEncode(mac.Sum(nil))

	return sigInput + "." + sig, nil
}

// NewLoginClaims 创建登录JWT claims
func NewLoginClaims(uin int64, deviceID string) *LoginClaims {
	now := time.Now().Unix()
	return &LoginClaims{
		Uin:        uin,
		Env:        0,
		Auth:       "web",
		TS:         now,
		APIID:      110,
		CltVersion: 79105,
		Src:        "man_machine.login_v3",
		DeviceID:   deviceID,
		ITS:        now,
		IAT:        now,
	}
}

// ParseJWTClaims 解析JWT载荷（不验签，仅解码）
func ParseJWTClaims(token string, out any) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return fmt.Errorf("JWT格式错误: 期望3段, 得到%d段", len(parts))
	}
	payload, err := base64URLDecode(parts[1])
	if err != nil {
		return fmt.Errorf("解码payload失败: %w", err)
	}
	return json.Unmarshal(payload, out)
}

func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func base64URLDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
