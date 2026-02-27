package main

import (
	"log"
	"miniprotocol/internal/auth"
	"miniprotocol/internal/friend"
	"miniprotocol/internal/httpc"
	"miniprotocol/internal/social"
)

// runFriendList 查询好友列表
func runFriendList(client *httpc.Client, cred *auth.Credential) {
	svc := friend.NewService(client, cred)
	resp, err := svc.List()
	if err != nil {
		log.Printf("[friend] 查询失败: %v", err)
		return
	}
	log.Printf("[friend] 共 %d 个好友:", resp.Data.Total)
	for _, f := range resp.Data.Friends {
		online := "离线"
		if f.IsOnline {
			online = "在线"
		}
		log.Printf("[friend]   %d %s (Lv%d) [%s]",
			f.Uin, f.NickName, f.Level, online)
	}
}

// runFriendAdd 发送好友申请
func runFriendAdd(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	svc := friend.NewService(client, cred)
	resp, err := svc.Add(targetUin, "")
	if err != nil {
		log.Printf("[friend] 添加失败: %v", err)
		return
	}
	log.Printf("[friend] 添加结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runFriendSearch 搜索用户
func runFriendSearch(client *httpc.Client, cred *auth.Credential, keyword string) {
	svc := friend.NewService(client, cred)
	resp, err := svc.Search(keyword)
	if err != nil {
		log.Printf("[friend] 搜索失败: %v", err)
		return
	}
	log.Printf("[friend] 搜索结果: %d 个", len(resp.Data.Users))
	for _, u := range resp.Data.Users {
		log.Printf("[friend]   %d %s (Lv%d)",
			u.Uin, u.NickName, u.Level)
	}
}

// runFriendDelete 删除好友
func runFriendDelete(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	svc := friend.NewService(client, cred)
	resp, err := svc.Delete(targetUin)
	if err != nil {
		log.Printf("[friend] 删除失败: %v", err)
		return
	}
	log.Printf("[friend] 删除结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runFriendAccept 接受好友申请
func runFriendAccept(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	svc := friend.NewService(client, cred)
	resp, err := svc.Accept(targetUin)
	if err != nil {
		log.Printf("[friend] 接受失败: %v", err)
		return
	}
	log.Printf("[friend] 接受结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runFriendRequests 查看好友申请列表
func runFriendRequests(client *httpc.Client, cred *auth.Credential) {
	svc := friend.NewService(client, cred)
	resp, err := svc.Requests()
	if err != nil {
		log.Printf("[friend] 查询申请列表失败: %v", err)
		return
	}
	log.Printf("[friend] 好友申请 %d 条:", len(resp.Data.Requests))
	for _, r := range resp.Data.Requests {
		log.Printf("[friend]   %d %s: %s", r.Uin, r.NickName, r.Message)
	}
}

// runOnlineFriends 查看在线好友
func runOnlineFriends(client *httpc.Client, cred *auth.Credential) {
	svc := friend.NewService(client, cred)
	resp, err := svc.Online()
	if err != nil {
		log.Printf("[friend] 查询在线好友失败: %v", err)
		return
	}
	log.Printf("[friend] 在线好友 %d 个:", len(resp.Data.Online))
	for _, uin := range resp.Data.Online {
		log.Printf("[friend]   %d", uin)
	}
}

// runFollow 关注用户
func runFollow(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	c := social.NewClient(client, cred)
	resp, err := c.Follow(targetUin)
	if err != nil {
		log.Printf("[social] 关注失败: %v", err)
		return
	}
	log.Printf("[social] 关注结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runUnfollow 取消关注
func runUnfollow(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	c := social.NewClient(client, cred)
	resp, err := c.Unfollow(targetUin)
	if err != nil {
		log.Printf("[social] 取消关注失败: %v", err)
		return
	}
	log.Printf("[social] 取消关注结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runLike 点赞用户
func runLike(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	c := social.NewClient(client, cred)
	resp, err := c.Like(targetUin, "user", "")
	if err != nil {
		log.Printf("[social] 点赞失败: %v", err)
		return
	}
	log.Printf("[social] 点赞结果: code=%d msg=%s", resp.Code, resp.Msg)
}

// runFansList 查询粉丝列表
func runFansList(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	c := social.NewClient(client, cred)
	resp, err := c.FansList(targetUin, 1, 20)
	if err != nil {
		log.Printf("[social] 粉丝列表失败: %v", err)
		return
	}
	log.Printf("[social] 粉丝共 %d 个:", resp.Data.Total)
	for _, u := range resp.Data.List {
		log.Printf("[social]   %d %s (Lv%d)", u.Uin, u.NickName, u.Level)
	}
}

// runFollowingList 查询关注列表
func runFollowingList(client *httpc.Client, cred *auth.Credential, targetUin int64) {
	c := social.NewClient(client, cred)
	resp, err := c.FollowList(targetUin, 1, 20)
	if err != nil {
		log.Printf("[social] 关注列表失败: %v", err)
		return
	}
	log.Printf("[social] 关注共 %d 个:", resp.Data.Total)
	for _, u := range resp.Data.List {
		log.Printf("[social]   %d %s (Lv%d)", u.Uin, u.NickName, u.Level)
	}
}
