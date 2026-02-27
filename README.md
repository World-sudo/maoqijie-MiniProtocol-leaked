# MiniProtocol

迷你世界 (MiniWorld) 协议的 Go 语言逆向实现。支持原生登录、SSO 登录、响应解密、聊天服务、游戏连接等完整流程。

## 功能

- **原生登录** — 逆向自 MicroMiniNew.exe，支持注册 + 密码登录 + Token 登录
- **SSO 登录** — Web 端登录流程，集成 GeeTest V4 验证码
- **响应解密** — 自定义 IV 重排 + AES-256-CBC 解密
- **RPC 通信** — XXTEA 加密 + DH 密钥交换
- **聊天服务** — IM 节点分配 + WebSocket 长连接 + RPC 调用
- **游戏连接** — WebSocket 二进制协议
- **遥测上报** — 事件采集与上报
- **辅助功能** — 防沉迷查询、信用分查询、内容审核、地图管理、版本检查等
- **好友系统** — 好友列表/搜索/添加/删除/接受申请
- **用户资料** — 查询资料/名片、修改昵称/签名/性别/皮肤
- **社交代理** — 关注/取关/点赞/粉丝列表/关注列表
- **邮件系统** — 邮件列表/读取/领取附件/删除
- **聊天消息解析** — IM消息帧结构解析(opcode+body)、序列化
- **游戏协议包** — 数据包解析/编码/分发处理框架
- **排行榜** — 全服榜/好友榜/地图榜/周榜
- **热更新** — 应用版本检查 + 增量补丁包查询

## 安装

```bash
go build -o mini.exe ./cmd/mini
```

依赖仅 `gorilla/websocket`，Go 1.22+。

## 使用

```bash
# 注册新账号
./mini -register -pwd <密码>

# 原生登录（仅登录）
./mini -uin <迷你号> -pwd <密码> -native -login-only

# SSO 登录
./mini -uin <迷你号> -pwd <密码>

# 完整连接（登录 + 聊天 + 游戏）
./mini -uin <迷你号> -jwt <token>

# 版本查询
./mini -version

# 内容审核
./mini -uin <迷你号> -jwt <token> -check-text <文本>

# 信用分查询
./mini -credit -uin <迷你号>

# 好友列表
./mini -uin <迷你号> -jwt <token> -friends

# 搜索用户
./mini -uin <迷你号> -jwt <token> -search <关键词>

# 查询用户资料
./mini -uin <迷你号> -jwt <token> -profile <目标uin>

# 关注用户
./mini -uin <迷你号> -jwt <token> -follow <目标uin>

# 邮件列表
./mini -uin <迷你号> -jwt <token> -mail

# 全服排行榜
./mini -uin <迷你号> -jwt <token> -rank

# 热更新检查
./mini -update
```

## 项目结构

```
cmd/
  mini/           主程序入口
  probe/          诊断探针，测试各模块功能

internal/
  register/       注册/登录（原生、SSO、旧版）
  auth/           签名、加密、JWT、凭证管理
  crypto/         XXTEA 算法实现
  rpc/            RPC 客户端 + DH 密钥交换
  chat/           聊天节点分配 + WebSocket 网关 + RPC
  game/           游戏 WebSocket 连接
  config/         全局配置与环境常量
  httpc/          HTTP 客户端 + MN-PAYLOAD 头部
  captcha/        GeeTest V4 验证码服务
  device/         设备指纹生成
  telemetry/      遥测事件上报
  room/           房间配置服务
  antiaddict/     防沉迷查询
  credit/         信用分查询
  moderation/     内容审核
  mapapi/         地图管理
  ugc/            资源上传
  version/        版本检查
  friend/         好友系统 (列表/搜索/添加/删除)
  profile/        用户资料 (查询/编辑)
  social/         社交代理 (关注/点赞/粉丝)
  mail/           邮件系统 (收发/领取)
  rank/           排行榜查询
  update/         热更新检查

```

## 协议概要

### 原生登录

```
POST https://wskacchm.mini1.cn:14130/login/auth_security
签名: md5("source=mini_micro&target=<t>&time=<ts>" + salt)
```

注册使用 `target=reg`，登录使用 `target=auth`，返回加密的 `authinfo` + `baseinfo`。

### 响应解密流程

```
iv 字段提取数字 N → 字符串循环右移 N 位 → Base64 解码 → AES-256-CBC 解密 → PKCS7 unpad
```

### RPC 通信

```
明文 → XXTEA 加密 → Base64 → URL 安全替换
密钥通过 DH 密钥交换协商（shared = B^a mod p → md5 分割为 Key + IV）
```

## 免责声明

本项目仅供学习和研究用途，请勿用于任何违反法律法规或游戏服务条款的行为。使用者需自行承担所有风险与责任。
