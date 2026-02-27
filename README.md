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
