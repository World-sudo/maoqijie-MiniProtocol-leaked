package auth

import (
	"fmt"
	"miniprotocol/internal/config"
	"time"
)

// Credential 聚合凭证管理
type Credential struct {
	Uin      int64
	DeviceID string
	LoginJWT string
	IMJWT    string
}

// NewCredential 创建凭证（需要已知 uin 和 deviceID）
func NewCredential(uin int64, deviceID string) *Credential {
	return &Credential{
		Uin:      uin,
		DeviceID: deviceID,
	}
}

// AuthString 生成认证字符串: switchAccountByAuthInfo_reg###<JWT>
func (c *Credential) AuthString() string {
	return config.AuthPrefix + c.LoginJWT
}

// ChatAuth 生成聊天auth签名
func (c *Credential) ChatAuth() string {
	return Sign(c.Uin, time.Now().Unix())
}

// ChatAuthAt 生成指定时间的聊天auth签名
func (c *Credential) ChatAuthAt(ts int64) string {
	return Sign(c.Uin, ts)
}

// SetLoginJWT 设置登录JWT
func (c *Credential) SetLoginJWT(token string) {
	c.LoginJWT = token
}

// SetIMJWT 设置IM JWT
func (c *Credential) SetIMJWT(token string) {
	c.IMJWT = token
}

// String 调试用
func (c *Credential) String() string {
	return fmt.Sprintf("Credential{uin=%d, device=%s, loginJWT=%d chars, imJWT=%d chars}",
		c.Uin, c.DeviceID, len(c.LoginJWT), len(c.IMJWT))
}
