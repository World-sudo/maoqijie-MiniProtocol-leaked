package version

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/url"
	"strconv"
)

// VersionInfo 版本信息响应
type VersionInfo struct {
	Version    string `json:"version"`
	CltVersion int    `json:"cltversion"`
	URL        string `json:"url"`
	ForceUp    int    `json:"force_up"`
	Desc       string `json:"desc"`
}

// PkgInfo 包信息
type PkgInfo struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	URL  string `json:"url"`
	MD5  string `json:"md5"`
	Size int64  `json:"size"`
}

// Checker 版本检查器
type Checker struct {
	client *httpc.Client
}

// NewChecker 创建版本检查器
func NewChecker(client *httpc.Client) *Checker {
	return &Checker{client: client}
}

// GetVersionJSON 获取版本信息
// GET https://mnweb.mini1.cn/game/version.json
func (c *Checker) GetVersionJSON() (*VersionInfo, error) {
	u := fmt.Sprintf("https://%s%s", config.VersionWebHost, config.VersionJSON)

	resp, err := c.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("获取版本信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取版本响应失败: %w", err)
	}

	var info VersionInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("解析版本信息失败: %w (body=%s)", err, string(body))
	}
	return &info, nil
}

// QueryPrimaryPkg 查询主包信息
// GET ?cmd=query_primary_pkg&apiid=110&cltversion=79106
func (c *Checker) QueryPrimaryPkg() (*PkgInfo, error) {
	return c.queryPkg("query_primary_pkg")
}

// QueryPatchPkg 查询补丁包信息
// GET ?cmd=query_patch_pkg&apiid=110&cltversion=79106
func (c *Checker) QueryPatchPkg() (*PkgInfo, error) {
	return c.queryPkg("query_patch_pkg")
}

func (c *Checker) queryPkg(cmd string) (*PkgInfo, error) {
	params := url.Values{}
	params.Set("cmd", cmd)
	params.Set("apiid", strconv.Itoa(config.APIID))
	params.Set("cltversion", strconv.Itoa(config.CltVersion))

	u := fmt.Sprintf("https://%s%s?%s",
		config.VersionWebHost, config.VersionJSON, params.Encode())

	resp, err := c.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("查询包信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取包信息失败: %w", err)
	}

	var info PkgInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("解析包信息失败: %w (body=%s)", err, string(body))
	}
	return &info, nil
}
