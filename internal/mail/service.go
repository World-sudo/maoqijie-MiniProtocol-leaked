// Package mail 邮件系统
// 逆向推断自 iworld.cfg: HWMailHost = "hwmail.mini1.cn" (海外)
// + shequ.mini1.cn 社区API 的 act= 接口模式
// + libiworld.dll 字符串: MailList, ReadMail, ReceiveMail, DeleteMail
//
// API端点推断:
//   shequ.mini1.cn:8080/miniw/mail/?act=<op>
//   mail.mini1.cn/mail/<op> (国内)
//   hwmail.mini1.cn/mail/<op> (海外)
//
// 签名方式与其他 shequ 接口一致
package mail

import (
	"encoding/json"
	"fmt"
	"io"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/config"
	"miniprotocol/internal/httpc"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// MailItem 邮件条目
// 逆向推断自 libiworld.dll: MailData 结构
type MailItem struct {
	MailID     string `json:"mail_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	SenderUin  int64  `json:"sender_uin"`
	SenderNick string `json:"sender_nick"`
	Type       int    `json:"type"` // 0=系统, 1=玩家, 2=活动奖励
	IsRead     bool   `json:"is_read"`
	HasAttach  bool   `json:"has_attach"`
	CreateTime int64  `json:"create_time"`
	ExpireTime int64  `json:"expire_time"`
}

// Attachment 邮件附件
type Attachment struct {
	ItemID   int    `json:"item_id"`
	ItemName string `json:"item_name"`
	Count    int    `json:"count"`
}

// ListResponse 邮件列表响应
type ListResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Total int        `json:"total"`
		Mails []MailItem `json:"mails"`
	} `json:"data"`
}

// DetailResponse 邮件详情响应
type DetailResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Mail    MailItem     `json:"mail"`
		Attachs []Attachment `json:"attachments"`
	} `json:"data"`
}

// ActionResponse 邮件操作响应
type ActionResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 邮件类型常量
const (
	TypeSystem  = 0 // 系统邮件
	TypePlayer  = 1 // 玩家邮件
	TypeReward  = 2 // 活动奖励
)

// Service 邮件服务
type Service struct {
	client *httpc.Client
	cred   *auth.Credential
}

// NewService 创建邮件服务
func NewService(client *httpc.Client, cred *auth.Credential) *Service {
	return &Service{client: client, cred: cred}
}

// List 获取邮件列表
// GET shequ.mini1.cn:8080/miniw/mail/?act=mail_list
func (s *Service) List() (*ListResponse, error) {
	params := s.baseParams()
	params.Set("act", "mail_list")
	return doGet[ListResponse](s.client, s.buildURL(params))
}

// ListByPage 分页获取邮件列表
func (s *Service) ListByPage(page, size int) (*ListResponse, error) {
	params := s.baseParams()
	params.Set("act", "mail_list")
	params.Set("page", strconv.Itoa(page))
	params.Set("size", strconv.Itoa(size))
	return doGet[ListResponse](s.client, s.buildURL(params))
}

// Read 读取邮件详情
// GET shequ.mini1.cn:8080/miniw/mail/?act=mail_read&mail_id=<id>
func (s *Service) Read(mailID string) (*DetailResponse, error) {
	params := s.baseParams()
	params.Set("act", "mail_read")
	params.Set("mail_id", mailID)
	return doGet[DetailResponse](s.client, s.buildURL(params))
}

// Receive 领取邮件附件
// POST shequ.mini1.cn:8080/miniw/mail/?act=mail_receive
func (s *Service) Receive(mailID string) (*ActionResponse, error) {
	form := s.baseForm()
	form.Set("mail_id", mailID)
	return s.doPost("mail_receive", form)
}

// ReceiveAll 一键领取所有附件
// POST shequ.mini1.cn:8080/miniw/mail/?act=mail_receive_all
func (s *Service) ReceiveAll() (*ActionResponse, error) {
	return s.doPost("mail_receive_all", s.baseForm())
}

// Delete 删除邮件
// POST shequ.mini1.cn:8080/miniw/mail/?act=mail_del
func (s *Service) Delete(mailID string) (*ActionResponse, error) {
	form := s.baseForm()
	form.Set("mail_id", mailID)
	return s.doPost("mail_del", form)
}

// DeleteRead 删除所有已读邮件
// POST shequ.mini1.cn:8080/miniw/mail/?act=mail_del_read
func (s *Service) DeleteRead() (*ActionResponse, error) {
	return s.doPost("mail_del_read", s.baseForm())
}

func (s *Service) buildURL(params url.Values) string {
	return fmt.Sprintf("http://%s:%d/miniw/mail/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())
}

func (s *Service) baseParams() url.Values {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(s.cred.Uin, 10)
	params := url.Values{}
	params.Set("uin", uinStr)
	params.Set("time", strconv.FormatInt(ts, 10))
	params.Set("auth", auth.Sign(s.cred.Uin, ts))
	params.Set("sign", auth.URLSignMD5(uinStr))
	params.Set("env", "0")
	return params
}

func (s *Service) baseForm() url.Values {
	return s.baseParams()
}

func (s *Service) doPost(act string, form url.Values) (*ActionResponse, error) {
	params := url.Values{}
	params.Set("act", act)

	u := fmt.Sprintf("http://%s:%d/miniw/mail/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	resp, err := s.client.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("邮件操作失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result ActionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
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
