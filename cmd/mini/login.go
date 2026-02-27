package main

import (
	"encoding/json"
	"fmt"
	"log"
	"miniprotocol/internal/captcha"
	"miniprotocol/internal/httpc"
	"miniprotocol/internal/register"
	"os"
)

// runRegister 注册新账号
func runRegister(client *httpc.Client, password, deviceID string) {
	native := register.NewNativeClient(client, "")

	log.Printf("[register] 注册中... device=%s", deviceID)

	resp, err := native.Register(password, deviceID)
	if err != nil {
		log.Fatalf("[register] 请求失败: %v", err)
	}

	log.Printf("[register] 响应: code=%d msg=%s", resp.Code, resp.Msg)

	if resp.Code != 0 {
		log.Fatalf("[register] 注册失败: code=%d msg=%s", resp.Code, resp.Msg)
	}

	log.Println("[register] 注册成功!")
	log.Printf("[register] iv: %s", resp.IV)

	// 解密 authinfo 获取 Uin/token/sign
	authData, err := resp.DecryptAuthInfo()
	if err != nil {
		log.Printf("[register] 解密authinfo失败: %v", err)
		log.Printf("[register] authinfo原文: %s...(%d字符)", resp.AuthInfo[:40], len(resp.AuthInfo))
	} else {
		log.Printf("[register] 解密成功! Uin=%d", authData.Uin)
		log.Printf("[register] token=%s", authData.Token)
		log.Printf("[register] sign=%s", authData.Sign)
	}

	// 解密 baseinfo
	baseData, err := resp.DecryptBaseInfo()
	if err != nil {
		log.Printf("[register] 解密baseinfo失败: %v", err)
	} else {
		log.Printf("[register] LastLoginTime=%d isLoginSafeVerify=%d",
			baseData.LastLoginTime, baseData.IsLoginSafeVerify)
	}

	// 保存凭证到文件
	cred := map[string]any{
		"authinfo": resp.AuthInfo,
		"baseinfo": resp.BaseInfo,
		"iv":       resp.IV,
		"deviceID": deviceID,
		"password": password,
	}
	if authData != nil {
		cred["uin"] = authData.Uin
		cred["token"] = authData.Token
		cred["sign"] = authData.Sign
	}
	data, _ := json.MarshalIndent(cred, "", "  ")
	filename := fmt.Sprintf("credential_%s.json", deviceID[3:11])
	if err := os.WriteFile(filename, data, 0600); err != nil {
		log.Printf("[register] 保存凭证失败: %v", err)
	} else {
		log.Printf("[register] 凭证已保存: %s", filename)
	}

	if authData != nil {
		log.Printf("[register] 你的迷你号(Uin): %d", authData.Uin)
		log.Printf("[register] 可用 -uin %d -pwd %s -native 登录", authData.Uin, password)
	}
}

// runNativeLogin 原生登录流程 (MicroMiniNew.exe 协议)
func runNativeLogin(client *httpc.Client, uin int64, password, deviceID string) (string, error) {
	native := register.NewNativeClient(client, "")

	log.Printf("[login] 原生登录: uin=%d device=%s", uin, deviceID)

	resp, err := native.Login(uin, password, deviceID)
	if err != nil {
		return "", fmt.Errorf("原生登录请求失败: %w", err)
	}

	log.Printf("[login] 响应: code=%d msg=%s", resp.Code, resp.Msg)

	if resp.Code != 0 {
		return "", fmt.Errorf("登录失败: code=%d msg=%s", resp.Code, resp.Msg)
	}

	// 打印登录信息
	if resp.Data != nil {
		log.Printf("[login] Uin=%d token=%s sign=%s",
			resp.Data.Uin, resp.Data.Token, resp.Data.Sign)
		return resp.Data.Token, nil
	}

	// data 字段为空时尝试从 authinfo 解密获取令牌
	if resp.AuthInfo != "" && resp.IV != "" {
		log.Printf("[login] data为空, 尝试解密authinfo...")
		authData, err := resp.DecryptAuthInfo()
		if err != nil {
			log.Printf("[login] 解密authinfo失败: %v", err)
		} else {
			log.Printf("[login] 解密成功! Uin=%d token=%s sign=%s",
				authData.Uin, authData.Token, authData.Sign)
			return authData.Token, nil
		}
	}

	return "", fmt.Errorf("登录响应缺少data字段且authinfo解密失败")
}

// runSSOLogin 执行SSO登录流程 (带手动验证码)
func runSSOLogin(client *httpc.Client, uin int64, password string) (string, error) {
	sso := register.NewSSOClient(client)

	log.Printf("[login] SSO登录: uin=%d", uin)
	deviceInfo := register.DefaultDevice()

	resp, err := sso.Login(uin, password, deviceInfo)
	if err != nil {
		return "", fmt.Errorf("SSO登录请求失败: %w", err)
	}

	log.Printf("[login] 响应: code=%d msg=%s", resp.Code, resp.Msg)

	if resp.Code == 0 {
		return extractToken(resp)
	}

	if !resp.NeedsCaptcha() {
		return "", fmt.Errorf("SSO返回异常: code=%d msg=%s", resp.Code, resp.Msg)
	}

	verID := resp.VerID()
	log.Printf("[login] 需要验证码, ver_id=%s", verID)

	captchaResult, err := captcha.Solve()
	if err != nil {
		return "", fmt.Errorf("验证码失败: %w", err)
	}

	log.Printf("[login] 验证码完成: lot=%s", captchaResult.LotNumber)

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

	switch v := resp.Data.(type) {
	case string:
		return v, nil
	case map[string]any:
		if token, ok := v["token"].(string); ok {
			return token, nil
		}
		raw, _ := json.Marshal(v)
		return "", fmt.Errorf("data中无token: %s", string(raw))
	default:
		raw, _ := json.Marshal(resp.Data)
		return "", fmt.Errorf("未知data类型: %s", string(raw))
	}
}
