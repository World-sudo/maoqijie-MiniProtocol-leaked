package register

// 响应解密模块
// 逆向自 MicroMiniNew.exe 响应解析函数 (0x0045D450) 和解密函数 (0x0045A180)
//
// 解密流程:
//  1. 从 iv 字段解析数字 N
//  2. 字符重排: rotated = encrypted[len-N:] + encrypted[:len-N]
//  3. 标准 Base64 解码
//  4. AES-256-CBC 解密 (Key/IV 初始化于 0x0064C440/0x0064C460)
//  5. PKCS7 unpad -> 明文 JSON

import (
	"encoding/json"
	"fmt"
	"miniprotocol/internal/auth"
)

// AuthInfoData authinfo 解密后的 JSON 结构
// 逆向自 MicroMiniNew.exe 0x0045D800-0x0045D950 字段解析代码
type AuthInfoData struct {
	Uin   int64  `json:"Uin"`
	Token string `json:"token"`
	Sign  string `json:"sign"`
}

// BaseInfoData baseinfo 解密后的 JSON 结构
// 已通过服务器验证: isloginsafeverify 实际为 bool 类型
type BaseInfoData struct {
	Uin               int64  `json:"Uin"`
	CreateTime        int64  `json:"CreateTime"`
	LastLoginTime     int64  `json:"LastLoginTime"`
	IsLoginSafeVerify bool   `json:"isloginsafeverify"`
	IsFreeze          bool   `json:"isfreeze"`
	Level             int    `json:"level"`
	APIID             int    `json:"apiid"`
	Phone             string `json:"Phone"`
	Email             string `json:"Email"`
	MiniCoin          int    `json:"minicoin"`
	MiniBean          int    `json:"minibean"`
	RoleInfo          *RoleInfo `json:"RoleInfo,omitempty"`
}

// RoleInfo 角色信息
type RoleInfo struct {
	Model      int    `json:"Model"`
	NickName   string `json:"NickName"`
	SkinID     int    `json:"SkinID"`
	SkinIDTime int    `json:"SkinIDTime"`
}

// decryptAuthInfo 通用 authinfo 解密
func decryptAuthInfo(authInfo, iv string) (*AuthInfoData, error) {
	if authInfo == "" || iv == "" {
		return nil, fmt.Errorf("authinfo or iv is empty")
	}
	plain, err := auth.DecryptNativeResponse(authInfo, iv)
	if err != nil {
		return nil, fmt.Errorf("decrypt authinfo: %w", err)
	}
	var data AuthInfoData
	if err := json.Unmarshal([]byte(plain), &data); err != nil {
		return nil, fmt.Errorf("parse authinfo JSON: %w (raw=%s)", err, plain)
	}
	return &data, nil
}

// decryptBaseInfo 通用 baseinfo 解密
func decryptBaseInfo(baseInfo, iv string) (*BaseInfoData, error) {
	if baseInfo == "" || iv == "" {
		return nil, fmt.Errorf("baseinfo or iv is empty")
	}
	plain, err := auth.DecryptNativeResponse(baseInfo, iv)
	if err != nil {
		return nil, fmt.Errorf("decrypt baseinfo: %w", err)
	}
	var data BaseInfoData
	if err := json.Unmarshal([]byte(plain), &data); err != nil {
		return nil, fmt.Errorf("parse baseinfo JSON: %w (raw=%s)", err, plain)
	}
	return &data, nil
}

// DecryptAuthInfo 解密注册响应的 authinfo 字段获取 Uin/token/sign
// 逆向自 MicroMiniNew.exe 0x0045D7CA -> 0x0045A180
func (r *RegisterResponse) DecryptAuthInfo() (*AuthInfoData, error) {
	return decryptAuthInfo(r.AuthInfo, r.IV)
}

// DecryptBaseInfo 解密注册响应的 baseinfo 字段
// 逆向自 MicroMiniNew.exe 0x0045D96F -> 0x0045A180
func (r *RegisterResponse) DecryptBaseInfo() (*BaseInfoData, error) {
	return decryptBaseInfo(r.BaseInfo, r.IV)
}

// DecryptAuthInfo 解密登录响应的 authinfo 字段
func (r *LoginResponse) DecryptAuthInfo() (*AuthInfoData, error) {
	return decryptAuthInfo(r.AuthInfo, r.IV)
}

// DecryptBaseInfo 解密登录响应的 baseinfo 字段
func (r *LoginResponse) DecryptBaseInfo() (*BaseInfoData, error) {
	return decryptBaseInfo(r.BaseInfo, r.IV)
}
