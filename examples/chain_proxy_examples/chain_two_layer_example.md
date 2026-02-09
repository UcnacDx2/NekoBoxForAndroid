# 基础双层代理配置示例 / Basic Two-Layer Proxy Configuration Example

[中文](#中文说明) | [English](#english-description)

---

## 中文说明

### 配置概述

这个示例展示如何配置一个基础的双层代理链：本地 SOCKS5 代理连接到远程 Shadowsocks 服务器。

### 使用场景

- 通过本地跳板访问远程服务器
- 隐藏真实 IP 地址
- 绕过本地网络限制

### 代理流程

```
客户端
  ↓
本地 SOCKS5 代理 (127.0.0.1:1080)
  ↓
远程 Shadowsocks 服务器 (example.com:8388)
  ↓
目标网站
```

## 配置步骤

### 第一步：添加本地 SOCKS5 代理

1. 在 NekoBox 主界面点击 "+" 按钮
2. 选择 "SOCKS5"
3. 填写配置信息：
   - **名称**: 本地SOCKS5
   - **服务器地址**: 127.0.0.1
   - **端口**: 1080
   - **用户名**: (如果需要)
   - **密码**: (如果需要)
4. 点击保存

### 第二步：添加 Shadowsocks 服务器

1. 点击 "+" 按钮
2. 选择 "Shadowsocks"
3. 填写配置信息：
   - **名称**: SS-US-Server
   - **服务器地址**: example.com
   - **端口**: 8388
   - **密码**: your-password
   - **加密方式**: aes-256-gcm
   - **插件**: (可选)
4. 点击保存

### 第三步：创建链式代理

1. 点击 "+" 按钮
2. 选择 "Proxy Chain" (链式代理)
3. 填写配置信息：
   - **名称**: 双层代理-本地-SS
4. 点击 "+" 添加代理：
   - 第一个：选择 "本地SOCKS5"
   - 第二个：选择 "SS-US-Server"
5. 确认代理顺序正确
6. 点击保存

### 第四步：使用链式代理

1. 返回主界面
2. 选择 "双层代理-本地-SS"
3. 点击连接按钮
4. 等待连接成功

## 配置参数详解

### 本地 SOCKS5 代理

| 参数 | 值 | 说明 |
|------|-----|------|
| 服务器地址 | 127.0.0.1 | 本地回环地址 |
| 端口 | 1080 | SOCKS5 默认端口 |
| 认证 | 可选 | 如果本地服务需要认证 |

### Shadowsocks 服务器

| 参数 | 值 | 说明 |
|------|-----|------|
| 服务器地址 | example.com | 你的服务器域名或 IP |
| 端口 | 8388 | 服务器端口 |
| 密码 | your-password | 服务器密码 |
| 加密方式 | aes-256-gcm | 推荐使用强加密 |

## 验证配置

### 1. 测试连接

在连接前，可以单独测试每个代理：

1. 先测试本地 SOCKS5 是否可用
2. 再测试 Shadowsocks 服务器是否可用
3. 最后测试链式代理

### 2. 检查日志

连接后，检查 NekoBox 日志：

1. 进入设置 → 日志
2. 查看连接日志
3. 确认没有错误信息

### 3. 验证 IP

访问 IP 查询网站，确认：
- 显示的 IP 是 Shadowsocks 服务器的 IP
- 不是你的真实 IP

## 故障排除

### 问题 1：无法连接到本地 SOCKS5

**可能原因**:
- 本地 SOCKS5 服务未启动
- 端口被占用
- 防火墙阻止

**解决方法**:
1. 确认本地 SOCKS5 服务正在运行
2. 检查端口 1080 是否可用：`netstat -an | grep 1080`
3. 检查防火墙设置

### 问题 2：无法连接到 Shadowsocks

**可能原因**:
- 服务器信息错误
- 服务器不可用
- 网络限制

**解决方法**:
1. 验证服务器地址、端口、密码
2. 单独测试 Shadowsocks 连接
3. 检查服务器状态

### 问题 3：速度很慢

**可能原因**:
- 服务器距离远
- 带宽限制
- 链式代理增加延迟

**解决方法**:
1. 选择地理位置更近的服务器
2. 检查每个代理的延迟
3. 考虑减少代理层数

## 优化建议

### 1. 性能优化
- 使用低延迟的本地代理
- 选择高速 Shadowsocks 服务器
- 启用 TCP Fast Open (如果支持)

### 2. 安全优化
- 使用强加密方式（aes-256-gcm）
- 启用混淆插件
- 定期更换密码

### 3. 稳定性优化
- 配置多个 Shadowsocks 备用服务器
- 使用可靠的本地代理服务
- 定期检查连接状态

---

## English Description

### Configuration Overview

This example demonstrates how to configure a basic two-layer proxy chain: local SOCKS5 proxy connecting to a remote Shadowsocks server.

### Use Cases

- Access remote servers through local relay
- Hide real IP address
- Bypass local network restrictions

### Proxy Flow

```
Client
  ↓
Local SOCKS5 Proxy (127.0.0.1:1080)
  ↓
Remote Shadowsocks Server (example.com:8388)
  ↓
Target Website
```

## Configuration Steps

### Step 1: Add Local SOCKS5 Proxy

1. Click "+" button on NekoBox main screen
2. Select "SOCKS5"
3. Fill in configuration:
   - **Name**: Local-SOCKS5
   - **Server Address**: 127.0.0.1
   - **Port**: 1080
   - **Username**: (if required)
   - **Password**: (if required)
4. Click save

### Step 2: Add Shadowsocks Server

1. Click "+" button
2. Select "Shadowsocks"
3. Fill in configuration:
   - **Name**: SS-US-Server
   - **Server Address**: example.com
   - **Port**: 8388
   - **Password**: your-password
   - **Encryption**: aes-256-gcm
   - **Plugin**: (optional)
4. Click save

### Step 3: Create Chain Proxy

1. Click "+" button
2. Select "Proxy Chain"
3. Fill in configuration:
   - **Name**: Two-Layer-Local-SS
4. Click "+" to add proxies:
   - First: Select "Local-SOCKS5"
   - Second: Select "SS-US-Server"
5. Confirm proxy order is correct
6. Click save

### Step 4: Use Chain Proxy

1. Return to main screen
2. Select "Two-Layer-Local-SS"
3. Click connect button
4. Wait for connection success

## Configuration Parameters

### Local SOCKS5 Proxy

| Parameter | Value | Description |
|-----------|-------|-------------|
| Server Address | 127.0.0.1 | Local loopback address |
| Port | 1080 | SOCKS5 default port |
| Authentication | Optional | If local service requires auth |

### Shadowsocks Server

| Parameter | Value | Description |
|-----------|-------|-------------|
| Server Address | example.com | Your server domain or IP |
| Port | 8388 | Server port |
| Password | your-password | Server password |
| Encryption | aes-256-gcm | Recommended strong encryption |

## Verify Configuration

### 1. Test Connection

Before connecting, test each proxy individually:

1. Test local SOCKS5 availability first
2. Then test Shadowsocks server availability
3. Finally test chain proxy

### 2. Check Logs

After connecting, check NekoBox logs:

1. Go to Settings → Logs
2. View connection logs
3. Confirm no error messages

### 3. Verify IP

Visit IP lookup website to confirm:
- Displayed IP is Shadowsocks server IP
- Not your real IP

## Troubleshooting

### Issue 1: Cannot Connect to Local SOCKS5

**Possible Causes**:
- Local SOCKS5 service not started
- Port occupied
- Firewall blocking

**Solutions**:
1. Confirm local SOCKS5 service is running
2. Check if port 1080 is available: `netstat -an | grep 1080`
3. Check firewall settings

### Issue 2: Cannot Connect to Shadowsocks

**Possible Causes**:
- Incorrect server information
- Server unavailable
- Network restrictions

**Solutions**:
1. Verify server address, port, password
2. Test Shadowsocks connection separately
3. Check server status

### Issue 3: Slow Speed

**Possible Causes**:
- Server far away
- Bandwidth limitation
- Chain proxy adds latency

**Solutions**:
1. Choose geographically closer servers
2. Check latency of each proxy
3. Consider reducing proxy layers

## Optimization Recommendations

### 1. Performance Optimization
- Use low-latency local proxy
- Choose high-speed Shadowsocks server
- Enable TCP Fast Open (if supported)

### 2. Security Optimization
- Use strong encryption (aes-256-gcm)
- Enable obfuscation plugin
- Regularly change passwords

### 3. Stability Optimization
- Configure multiple backup Shadowsocks servers
- Use reliable local proxy service
- Regularly check connection status

---

## 许可证 / License

GPL-3.0
