// Package update 热更新检查
// 逆向自 libMiniBaseEngine.dll 热更新模块:
//   update.mini1.cn:13002 - 补丁服务器
//   /miniw/patch_server - 补丁包查询路径
//   mwu-api-pre.mini1.cn/app_update/check_app_ver - 应用版本检查
//
// 热更新流程:
//   1. 客户端携带当前 cltversion 请求检查
//   2. 服务端返回最新版本信息 + 补丁包URL
//   3. 客户端下载补丁包 (增量/全量)
//   4. 校验MD5 → 解压替换 → 重启游戏
//
// 参数签名: token = auth.Sign(uin, ts), uin + channel + env + os_type + app_ver
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/url"
	"strconv"
	"time"
)

// CheckResponse 版本检查响应
type CheckResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		NeedUpdate  bool   `json:"need_update"`
		LatestVer   int    `json:"latest_ver"`
		DownloadURL string `json:"download_url"`
		MD5         string `json:"md5"`
		Size        int64  `json:"size"`
		ForceUpdate bool   `json:"force_update"`
		Desc        string `json:"desc"`
	} `json:"data"`
}

// PatchResponse 补丁包查询响应
type PatchResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Patches []PatchInfo `json:"patches"`
	} `json:"data"`
}

// PatchInfo 补丁包信息
type PatchInfo struct {
	FromVer int    `json:"from_ver"`
	ToVer   int    `json:"to_ver"`
	URL     string `json:"url"`
	MD5     string `json:"md5"`
	Size    int64  `json:"size"`
}

// Checker 热更新检查器
type Checker struct {
	client *httpc.Client
}

// NewChecker 创建热更新检查器
func NewChecker(client *httpc.Client) *Checker {
	return &Checker{client: client}
}

// CheckAppVersion 检查应用版本更新
// GET mwu-api-pre.mini1.cn/app_update/check_app_ver
// 逆向自 register.BuildUpdateCheckURL
func (c *Checker) CheckAppVersion(uin int64) (*CheckResponse, error) {
	ts := time.Now().Unix()
	token := auth.Sign(uin, ts)

	params := url.Values{}
	params.Set("app_ver", strconv.Itoa(config.CltVersion))
	params.Set("channel", strconv.Itoa(config.APIID))
	params.Set("env", "0")
	params.Set("os_type", "1") // 1=Windows
	params.Set("token", token)
	params.Set("uin", strconv.FormatInt(uin, 10))

	u := fmt.Sprintf("https://mwu-api-pre.mini1.cn/app_update/check_app_ver?%s",
		params.Encode())

	return doGet[CheckResponse](c.client, u)
}

// QueryPatches 查询增量补丁包
// GET update.mini1.cn:13002/miniw/patch_server
// 逆向自 config.UpdateServer + config.UpdatePkgPath
func (c *Checker) QueryPatches(currentVer int) (*PatchResponse, error) {
	params := url.Values{}
	params.Set("cmd", "query_patch")
	params.Set("apiid", strconv.Itoa(config.APIID))
	params.Set("cltversion", strconv.Itoa(currentVer))
	params.Set("os_type", "1")
	params.Set("env", "0")

	u := fmt.Sprintf("https://%s%s?%s",
		config.UpdateServer, config.UpdatePkgPath, params.Encode())

	return doGet[PatchResponse](c.client, u)
}

// QueryLatestPatches 使用当前版本查询补丁
func (c *Checker) QueryLatestPatches() (*PatchResponse, error) {
	return c.QueryPatches(config.CltVersion)
}

// CheckEngine 检查引擎资源更新
// GET engine.mini1.cn/<path>
func (c *Checker) CheckEngine(path string) ([]byte, error) {
	u := fmt.Sprintf("https://%s/%s", config.EngineAssets, path)

	resp, err := c.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("引擎资源检查失败: %w", err)
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func doGet[T any](client *httpc.Client, u string) (*T, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
