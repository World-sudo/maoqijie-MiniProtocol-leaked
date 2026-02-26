// Package captcha 本地验证码服务
// 启动 HTTP 服务器展示 GeeTest V4 验证码页面
// 用户在浏览器中手动完成验证，结果回传给 Go 程序
package captcha

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"miniprotocol/internal/config"
	"net"
	"net/http"
	"time"
)

// Result GeeTest V4 验证结果
// 逆向自 sso.mini1.cn JS: captchaInstance.getValidate()
type Result struct {
	CaptchaID     string `json:"captcha_id"`
	LotNumber     string `json:"lot_number"`
	CaptchaOutput string `json:"captcha_output"`
	PassToken     string `json:"pass_token"`
	GenTime       string `json:"gen_time"`
}

// Solve 启动本地服务器让用户手动完成验证码
// 返回验证结果，超时 5 分钟自动取消
func Solve() (*Result, error) {
	return SolveWithID(config.GeeTestCaptchaID)
}

// SolveWithID 使用指定 captchaID 完成验证
func SolveWithID(captchaID string) (*Result, error) {
	resultCh := make(chan *Result, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		fmt.Fprintf(w, captchaPage, captchaID)
	})
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", 405)
			return
		}
		var result Result
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, "bad request", 400)
			errCh <- fmt.Errorf("解析验证结果失败: %w", err)
			return
		}
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.Write([]byte(successPage))
		resultCh <- &result
	})

	// 动态选择可用端口
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("监听端口失败: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	srv := &http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	log.Printf("[captcha] 请在浏览器打开: %s", url)

	// 等待结果或超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	defer srv.Shutdown(context.Background())

	select {
	case result := <-resultCh:
		log.Println("[captcha] 验证码完成")
		return result, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, fmt.Errorf("验证码超时 (5分钟)")
	}
}
