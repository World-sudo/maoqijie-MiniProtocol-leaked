package main

import (
	"encoding/json"
	"fmt"
	"log"
	"miniprotocol/internal/captcha"
	"miniprotocol/internal/httpc"
	"miniprotocol/internal/register"
)

// runNativeLogin 原生登录流程 (MicroMiniNew.exe 协议)
// API: POST https://wskacchm.mini1.cn:14130/login/auth_security
func runNativeLogin(client *httpc.Client, uin int64, password, deviceID string) (string, error) {
	native := register.NewNativeClient(client, "")

	log.Printf("[login] 原生登录: uin=%d device=%s", uin, deviceID)

	resp, err := native.Login(uin, password, deviceID)
	if err != nil {
		return "", fmt.Errorf("原生登录请求失败: %w", err)
	}

	log.Printf("[login] 响应: code=%d msg=%s", resp.Code, resp.Msg)

	if resp.Code != 0 {
		raw, _ := json.Marshal(resp)
		return "", fmt.Errorf("原生登录失败: %s", string(raw))
	}

	// 解析登录成功数据
	dataBytes, _ := json.Marshal(resp.Data)
	var loginData register.LoginData
	if err := json.Unmarshal(dataBytes, &loginData); err != nil {
		return "", fmt.Errorf("解析登录数据失败: %w", err)
	}

	log.Printf("[login] Uin=%d token=%s", loginData.Uin, loginData.Token)
	return loginData.Token, nil
}

// runSSOLogin 执行SSO登录流程 (带手动验证码)
// 返回: uin, jwt token, error
func runSSOLogin(client *httpc.Client, uin int64, password string) (string, error) {
	sso := register.NewSSOClient(client)

	// 第一步: 发送登录请求 (不带验证码)
	log.Printf("[login] SSO登录: uin=%d", uin)
	deviceInfo := register.DefaultDevice()

	resp, err := sso.Login(uin, password, deviceInfo)
	if err != nil {
		return "", fmt.Errorf("SSO登录请求失败: %w", err)
	}

	log.Printf("[login] 响应: code=%d msg=%s", resp.Code, resp.Msg)

	// 直接成功 (罕见，通常需要验证码)
	if resp.Code == 0 {
		return extractToken(resp)
	}

	// 需要验证码
	if !resp.NeedsCaptcha() {
		return "", fmt.Errorf("SSO返回异常: code=%d msg=%s", resp.Code, resp.Msg)
	}

	verID := resp.VerID()
	log.Printf("[login] 需要验证码, ver_id=%s", verID)

	// 第二步: 弹出验证码让用户手动完成
	captchaResult, err := captcha.Solve()
	if err != nil {
		return "", fmt.Errorf("验证码失败: %w", err)
	}

	log.Printf("[login] 验证码完成: lot=%s", captchaResult.LotNumber)

	// 第三步: 带验证码重新登录
	gt := &register.GeeTestData{
		Platform:      "web-sso",
		Version:       "",
		CaptchaID:     captchaResult.CaptchaID,
		LotNumber:     captchaResult.LotNumber,
		CaptchaOutput: captchaResult.CaptchaOutput,
		PassToken:     captchaResult.PassToken,
		GenTime:       captchaResult.GenTime,
	}

	resp2, err := sso.LoginWithCaptcha(uin, password, deviceInfo, verID, gt)
	if err != nil {
		return "", fmt.Errorf("SSO验证码登录失败: %w", err)
	}

	log.Printf("[login] 验证码登录响应: code=%d msg=%s", resp2.Code, resp2.Msg)

	if resp2.Code != 0 {
		return "", fmt.Errorf("登录失败: code=%d msg=%s", resp2.Code, resp2.Msg)
	}

	return extractToken(resp2)
}

// extractToken 从SSO响应中提取JWT令牌
func extractToken(resp *register.SSOLoginResponse) (string, error) {
	if resp.Data == nil {
		return "", fmt.Errorf("响应无data字段")
	}

	// data 可能是字符串(token)或对象
	switch v := resp.Data.(type) {
	case string:
		return v, nil
	case map[string]any:
		if token, ok := v["token"].(string); ok {
			return token, nil
		}
		// 打印完整data便于调试
		raw, _ := json.Marshal(v)
		return "", fmt.Errorf("data中无token: %s", string(raw))
	default:
		raw, _ := json.Marshal(resp.Data)
		return "", fmt.Errorf("未知data类型: %s", string(raw))
	}
}
