// Package rank 排行榜查询
// 逆向推断自 shequ.mini1.cn 社区API act= 接口模式
// + libiworld.dll 字符串: RankList, FriendRank, GlobalRank, MapRank
//
// API端点推断:
//   shequ.mini1.cn:8080/miniw/rank/?act=<rank_type>
//
// 排行榜类型:
//   rank_global - 全服排行 (等级/经验)
//   rank_friend - 好友排行
//   rank_map    - 地图热度排行
//   rank_weekly - 周榜
package rank

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

// RankItem 排行榜条目
type RankItem struct {
	Rank     int    `json:"rank"`
	Uin      int64  `json:"uin"`
	NickName string `json:"nick_name"`
	Level    int    `json:"level"`
	Score    int64  `json:"score"`
	VIP      int    `json:"vip"`
	SkinID   int    `json:"skin_id"`
}

// RankResponse 排行榜响应
type RankResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total  int        `json:"total"`
		MyRank int        `json:"my_rank"`
		List   []RankItem `json:"list"`
	} `json:"data"`
}

// MapRankItem 地图排行榜条目
type MapRankItem struct {
	Rank     int    `json:"rank"`
	MapID    string `json:"map_id"`
	MapName  string `json:"map_name"`
	Author   string `json:"author"`
	AuthUin  int64  `json:"author_uin"`
	PlayCnt  int64  `json:"play_count"`
	LikeCnt  int64  `json:"like_count"`
}

// MapRankResponse 地图排行榜响应
type MapRankResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total int           `json:"total"`
		List  []MapRankItem `json:"list"`
	} `json:"data"`
}

// 排行榜类型常量
const (
	TypeGlobal = "rank_global" // 全服等级榜
	TypeFriend = "rank_friend" // 好友排行
	TypeMap    = "rank_map"    // 地图热度榜
	TypeWeekly = "rank_weekly" // 周榜
)

// Client 排行榜查询客户端
type Client struct {
	httpClient *httpc.Client
	cred       *auth.Credential
}

// NewClient 创建排行榜客户端
func NewClient(c *httpc.Client, cred *auth.Credential) *Client {
	return &Client{httpClient: c, cred: cred}
}

// Global 全服排行榜
// GET shequ.mini1.cn:8080/miniw/rank/?act=rank_global&page=<p>&size=<s>
func (c *Client) Global(page, size int) (*RankResponse, error) {
	return c.queryRank(TypeGlobal, page, size)
}

// Friend 好友排行榜
func (c *Client) Friend(page, size int) (*RankResponse, error) {
	return c.queryRank(TypeFriend, page, size)
}

// Weekly 周榜
func (c *Client) Weekly(page, size int) (*RankResponse, error) {
	return c.queryRank(TypeWeekly, page, size)
}

// Map 地图热度排行
// GET shequ.mini1.cn:8080/miniw/rank/?act=rank_map&page=<p>&size=<s>
func (c *Client) Map(page, size int) (*MapRankResponse, error) {
	params := c.baseParams()
	params.Set("act", TypeMap)
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	u := fmt.Sprintf("http://%s:%d/miniw/rank/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[MapRankResponse](c.httpClient, u)
}

func (c *Client) queryRank(rankType string, page, size int) (*RankResponse, error) {
	params := c.baseParams()
	params.Set("act", rankType)
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))

	u := fmt.Sprintf("http://%s:%d/miniw/rank/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	return doGet[RankResponse](c.httpClient, u)
}

func (c *Client) baseParams() url.Values {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	params := url.Values{}
	params.Set("uin", uinStr)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("auth", auth.Sign(c.cred.Uin, ts))
	params.Set("sign", auth.URLSignMD5(uinStr))
	params.Set("env", "0")
	params.Set("country", "CN")
	params.Set("lang", "1")
	return params
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
