package main

import (
	"log"
	"miniprotocol/internal/antiaddict"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/httpc"
	"miniprotocol/internal/mail"
	"miniprotocol/internal/profile"
	"miniprotocol/internal/rank"
	"miniprotocol/internal/update"
)

// runProfileQuery 查询用户资料
func runProfileQuery(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	c := profile.NewClient(client, cred)
	resp, err := c.GetInfo(targetUin)
	if err != nil {
		log.Printf("[profile] 查询失败: %v", err)
		return
	}
	if resp.Data == nil {
		log.Printf("[profile] code=%d msg=%s", resp.Code, resp.Msg)
		return
	}
	d := resp.Data
	log.Printf("[profile] %d %s (Lv%d VIP%d)", d.Uin, d.NickName, d.Level, d.VIP)
	log.Printf("[profile] 签名: %s", d.Sign)
	log.Printf("[profile] 粉丝=%d 关注=%d 好友=%d 地图=%d",
		d.FansCount, d.FollowCnt, d.FriendCnt, d.MapCount)
}

// runProfileCard 查询用户名片
func runProfileCard(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	c := profile.NewClient(client, cred)
	resp, err := c.GetCard(targetUin)
	if err != nil {
		log.Printf("[profile] 名片查询失败: %v", err)
		return
	}
	if resp.Data == nil {
		log.Printf("[profile] code=%d msg=%s", resp.Code, resp.Msg)
		return
	}
	d := resp.Data
	log.Printf("[profile] 名片: %d %s (Lv%d VIP%d) 签名=%s",
		d.Uin, d.NickName, d.Level, d.VIP, d.Sign)
}

// runSetNick 修改昵称
func runSetNick(client *httpc.Client, cred *auth.Credential, nick string) {
	c := profile.NewClient(client, cred)
	resp, err := c.SetNickName(nick)
	if err != nil {
		log.Printf("[profile] 修改昵称失败: %v", err)
		return
	}
	log.Printf("[profile] 修改昵称: code=%d msg=%s", resp.Code, resp.Msg)
}

// runSetSign 修改签名
func runSetSign(client *httpc.Client, cred *auth.Credential, sign string) {
	c := profile.NewClient(client, cred)
	resp, err := c.SetSign(sign)
	if err != nil {
		log.Printf("[profile] 修改签名失败: %v", err)
		return
	}
	log.Printf("[profile] 修改签名: code=%d msg=%s", resp.Code, resp.Msg)
}

// runSetGender 修改性别
func runSetGender(client *httpc.Client, cred *auth.Credential, gender int) {
	c := profile.NewClient(client, cred)
	resp, err := c.SetGender(gender)
	if err != nil {
		log.Printf("[profile] 修改性别失败: %v", err)
		return
	}
	log.Printf("[profile] 修改性别: code=%d msg=%s", resp.Code, resp.Msg)
}

// runSetSkin 修改皮肤
func runSetSkin(client *httpc.Client, cred *auth.Credential, skinID int) {
	c := profile.NewClient(client, cred)
	resp, err := c.SetSkin(skinID)
	if err != nil {
		log.Printf("[profile] 修改皮肤失败: %v", err)
		return
	}
	log.Printf("[profile] 修改皮肤: code=%d msg=%s", resp.Code, resp.Msg)
}

// runMailList 查询邮件列表
func runMailList(client *httpc.Client, cred *auth.Credential) {
	svc := mail.NewService(client, cred)
	resp, err := svc.List()
	if err != nil {
		log.Printf("[mail] 查询失败: %v", err)
		return
	}
	log.Printf("[mail] 共 %d 封邮件:", resp.Data.Total)
	for _, m := range resp.Data.Mails {
		read := "未读"
		if m.IsRead {
			read = "已读"
		}
		log.Printf("[mail]   [%s] %s - %s (%s)",
			m.MailID, m.Title, m.SenderNick, read)
	}
}

// runMailReceiveAll 一键领取所有邮件附件
func runMailReceiveAll(client *httpc.Client, cred *auth.Credential) {
	svc := mail.NewService(client, cred)
	resp, err := svc.ReceiveAll()
	if err != nil {
		log.Printf("[mail] 领取失败: %v", err)
		return
	}
	log.Printf("[mail] 领取结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runMailDelete 删除指定邮件
func runMailDelete(client *httpc.Client, cred *auth.Credential, mailID string) {
	svc := mail.NewService(client, cred)
	resp, err := svc.Delete(mailID)
	if err != nil {
		log.Printf("[mail] 删除失败: %v", err)
		return
	}
	log.Printf("[mail] 删除结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runRankGlobal 全服排行榜
func runRankGlobal(client *httpc.Client, cred *auth.Credential) {
	c := rank.NewClient(client, cred)
	resp, err := c.Global(1, 20)
	if err != nil {
		log.Printf("[rank] 查询失败: %v", err)
		return
	}
	log.Printf("[rank] 全服排行 (我的排名: %d):", resp.Data.MyRank)
	for _, r := range resp.Data.List {
		log.Printf("[rank]   #%d %d %s Lv%d 分数=%d",
			r.Rank, r.Uin, r.NickName, r.Level, r.Score)
	}
}

// runRankFriend 好友排行榜
func runRankFriend(client *httpc.Client, cred *auth.Credential) {
	c := rank.NewClient(client, cred)
	resp, err := c.Friend(1, 20)
	if err != nil {
		log.Printf("[rank] 查询失败: %v", err)
		return
	}
	log.Printf("[rank] 好友排行 (我的排名: %d):", resp.Data.MyRank)
	for _, r := range resp.Data.List {
		log.Printf("[rank]   #%d %d %s Lv%d 分数=%d",
			r.Rank, r.Uin, r.NickName, r.Level, r.Score)
	}
}

// runRankWeekly 周榜
func runRankWeekly(client *httpc.Client, cred *auth.Credential) {
	c := rank.NewClient(client, cred)
	resp, err := c.Weekly(1, 20)
	if err != nil {
		log.Printf("[rank] 查询失败: %v", err)
		return
	}
	log.Printf("[rank] 周榜 (我的排名: %d):", resp.Data.MyRank)
	for _, r := range resp.Data.List {
		log.Printf("[rank]   #%d %d %s Lv%d 分数=%d",
			r.Rank, r.Uin, r.NickName, r.Level, r.Score)
	}
}

// runRankMap 地图排行榜
func runRankMap(client *httpc.Client, cred *auth.Credential) {
	c := rank.NewClient(client, cred)
	resp, err := c.Map(1, 20)
	if err != nil {
		log.Printf("[rank] 查询失败: %v", err)
		return
	}
	log.Printf("[rank] 地图排行 共 %d:", resp.Data.Total)
	for _, r := range resp.Data.List {
		log.Printf("[rank]   #%d %s by %s 玩=%d 赞=%d",
			r.Rank, r.MapName, r.Author, r.PlayCnt, r.LikeCnt)
	}
}

// runAntiAddict 防沉迷查询
func runAntiAddict(client *httpc.Client, uin int64) {
	c := antiaddict.NewClient(client)
	resp, err := c.Query(uin)
	if err != nil {
		log.Printf("[antiaddict] 查询失败: %v", err)
		return
	}
	log.Printf("[antiaddict] code=%d 状态=%d 时长=%d 剩余=%d",
		resp.Code, resp.Data.Status, resp.Data.Duration, resp.Data.Remain)
}

// runUpdateCheck 检查热更新
func runUpdateCheck(client *httpc.Client, uin int64) {
	checker := update.NewChecker(client)

	resp, err := checker.CheckAppVersion(uin)
	if err != nil {
		log.Printf("[update] 版本检查失败: %v", err)
		return
	}
	if resp.Data.NeedUpdate {
		log.Printf("[update] 有新版本: v%d", resp.Data.LatestVer)
		log.Printf("[update] 下载: %s (大小=%d)", resp.Data.DownloadURL, resp.Data.Size)
		if resp.Data.ForceUpdate {
			log.Printf("[update] 强制更新: %s", resp.Data.Desc)
		}
	} else {
		log.Println("[update] 当前已是最新版本")
	}

	patchResp, err := checker.QueryLatestPatches()
	if err != nil {
		log.Printf("[update] 补丁查询失败: %v", err)
		return
	}
	for _, p := range patchResp.Data.Patches {
		log.Printf("[update] 补丁: v%d→v%d %s (大小=%d)",
			p.FromVer, p.ToVer, p.URL, p.Size)
	}
}
