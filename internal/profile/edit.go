package profile

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

// EditResponse 资料编辑响应
type EditResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// SetNickName 修改昵称
// POST shequ.mini1.cn:8080/miniw/user/?act=set_nick
// 逆向推断自 libiworld.dll: changeNickName / modifyUserNick
func (c *Client) SetNickName(nick string) (*EditResponse, error) {
	form := c.baseForm()
	form.Set("nick_name", nick)
	return c.doPost("set_nick", form)
}

// SetSign 修改个性签名
// POST shequ.mini1.cn:8080/miniw/user/?act=set_sign
func (c *Client) SetSign(sign string) (*EditResponse, error) {
	form := c.baseForm()
	form.Set("sign", sign)
	return c.doPost("set_sign", form)
}

// SetGender 修改性别 (0=未知, 1=男, 2=女)
// POST shequ.mini1.cn:8080/miniw/user/?act=set_gender
func (c *Client) SetGender(gender int) (*EditResponse, error) {
	form := c.baseForm()
	form.Set("gender", strconv.Itoa(gender))
	return c.doPost("set_gender", form)
}

// SetBirthday 修改生日
// POST shequ.mini1.cn:8080/miniw/user/?act=set_birthday
func (c *Client) SetBirthday(birthday string) (*EditResponse, error) {
	form := c.baseForm()
	form.Set("birthday", birthday)
	return c.doPost("set_birthday", form)
}

// SetSkin 修改皮肤
// POST shequ.mini1.cn:8080/miniw/user/?act=set_skin
// 逆向推断自 libiworld.dll: SetSkinID
func (c *Client) SetSkin(skinID int) (*EditResponse, error) {
	form := c.baseForm()
	form.Set("skin_id", strconv.Itoa(skinID))
	return c.doPost("set_skin", form)
}

func (c *Client) baseForm() url.Values {
	ts := time.Now().Unix()
	uinStr := strconv.FormatInt(c.cred.Uin, 10)
	form := url.Values{}
	form.Set("uin", uinStr)
	form.Set("time", strconv.FormatInt(ts, 10))
	form.Set("auth", auth.Sign(c.cred.Uin, ts))
	form.Set("sign", auth.URLSignMD5(uinStr))
	form.Set("env", "0")
	form.Set("token", c.cred.LoginJWT)
	return form
}

func (c *Client) doPost(act string, form url.Values) (*EditResponse, error) {
	params := url.Values{}
	params.Set("act", act)

	u := fmt.Sprintf("http://%s:%d/miniw/user/?%s",
		config.Servers[config.EnvDomestic].ShequHTTP,
		config.ShequHTTPPort, params.Encode())

	resp, err := c.httpClient.Post(u, "application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("资料编辑请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result EditResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w (body=%s)", err, string(body))
	}
	return &result, nil
}
