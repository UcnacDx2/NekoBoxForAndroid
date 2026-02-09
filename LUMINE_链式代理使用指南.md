# 如何在NekoBoxForAndroid中使用Lumine链式代理

## 概述

本指南介绍如何配置NekoBoxForAndroid，使用lumine作为上游代理（链式代理），为您的VPN连接提供TLS分片和审查绕过功能。

## 什么是链式代理？

链式代理（也称为上游代理或代理链）是指按顺序通过多个代理路由VPN流量：

```
您的设备 → NekoBoxForAndroid → Lumine（预处理） → 您的代理服务器 → 互联网
```

优势：
- **TLS分片**: Lumine在到达代理服务器之前对TLS握手进行分片
- **绕过DPI**: 帮助在本地网络级别绕过深度包检测
- **基于策略**: 针对不同域名/IP使用不同的分片策略
- **透明**: 适用于任何代理协议（SOCKS5, HTTP, Shadowsocks, VMess等）

## 快速开始

### 方法1: 使用Lumine配置文件

1. **创建lumine配置文件** (例如 `/sdcard/lumine_config.json`):

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "8.8.8.8:53",
  "udp_minsize": 4096,
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  },
  "domain_policies": {},
  "ip_policies": {}
}
```

2. **在NekoBoxForAndroid配置中**，添加lumine作为SOCKS5出站：

```json
{
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080,
      "version": "5"
    },
    {
      "type": "shadowsocks",
      "tag": "proxy",
      "server": "your.server.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "your-password",
      "detour": "lumine"
    }
  ]
}
```

其中 `"detour": "lumine"` 使代理首先通过lumine路由。

### 方法2: 编程方式配置

在VPN服务中添加以下代码：

```kotlin
import libcore.Libcore
import libcore.LumineService
import java.io.File

class MyVpnService {
    private var lumineService: LumineService? = null
    
    fun startLumine() {
        // 1. 启动lumine服务
        val configFile = File("/sdcard/lumine_config.json")
        lumineService = if (configFile.exists()) {
            Libcore.newLumineService(configFile.absolutePath)
        } else {
            // 使用默认地址
            Libcore.newLumineServiceWithAddrs(
                "127.0.0.1:1080",  // SOCKS5地址
                "none",            // HTTP地址（不启用）
                ""                 // 配置JSON（可选）
            )
        }
        
        // 2. 配置sing-box使用lumine作为detour
        // 在您的sing-box配置中添加lumine出站和detour配置
    }
    
    fun stopLumine() {
        lumineService?.close()
        lumineService = null
    }
}
```

## 配置示例

### 示例1: Shadowsocks with Lumine预处理

完整的sing-box配置文件 (`singbox_lumine_chain_example.json`):

```json
{
  "log": {"level": "info"},
  "dns": {
    "servers": [
      {
        "tag": "dns_proxy",
        "address": "tls://8.8.8.8",
        "detour": "proxy"
      }
    ]
  },
  "inbounds": [
    {
      "type": "tun",
      "tag": "tun-in",
      "inet4_address": "172.19.0.1/28",
      "auto_route": true,
      "sniff": true
    }
  ],
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080,
      "version": "5"
    },
    {
      "type": "shadowsocks",
      "tag": "proxy",
      "server": "example.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "password",
      "detour": "lumine"
    },
    {
      "type": "direct",
      "tag": "direct"
    }
  ],
  "route": {
    "auto_detect_interface": true,
    "rules": [
      {"geoip": "cn", "outbound": "direct"}
    ],
    "final": "proxy"
  }
}
```

### 示例2: VMess with Lumine预处理

参见 `singbox_vmess_lumine_example.json` 文件。

## Lumine配置选项

### 基础配置

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "8.8.8.8:53"
}
```

### TLS分片高级配置

```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  }
}
```

- `mode`: 模式选择
  - `"tls-rf"`: TLS记录分片
  - `"ttl-d"`: TTL去同步
  - `"direct"`: 不处理
- `num_records`: 分片记录数（越高=越多分片）
- `num_segs`: 每个记录的TCP段数
- `send_interval`: 片段之间的延迟

### 针对特定域名的策略

```json
{
  "domain_policies": {
    "*.google.com;*.youtube.com": {
      "mode": "tls-rf",
      "num_records": 15,
      "num_segs": 4,
      "send_interval": "250ms"
    },
    "*.github.com": {
      "mode": "tls-rf",
      "num_records": 8,
      "num_segs": 2,
      "send_interval": "100ms"
    }
  }
}
```

参见 `lumine_config_advanced.json` 获取更多示例。

## 使用步骤

### 1. 准备配置文件

将配置文件放在可访问的位置，例如：
- `/sdcard/Download/lumine_config.json`
- `/sdcard/Documents/lumine_config.json`
- 应用的私有目录

### 2. 启动Lumine服务

在VPN启动时，先启动lumine：

```kotlin
// 在VPN Service的onCreate()或onStartCommand()中
val configPath = "/sdcard/lumine_config.json"
lumineService = Libcore.newLumineService(configPath)
```

### 3. 配置sing-box

在sing-box配置的outbounds中：

1. 添加lumine SOCKS5出站
2. 在主代理outbound中添加 `"detour": "lumine"`

### 4. 停止Lumine服务

在VPN停止时，关闭lumine：

```kotlin
// 在VPN Service的onDestroy()中
lumineService?.close()
lumineService = null
```

## 故障排除

### Lumine无法启动

检查端口是否可用：
```bash
netstat -an | grep 1080
```

尝试使用不同端口：
```json
{
  "socks5_address": "127.0.0.1:10800"
}
```

### 连接超时

- 检查lumine是否在运行
- 验证SOCKS5地址与lumine配置匹配
- 检查防火墙/SELinux设置

### 没有分片效果

- 验证 `mode` 设置为 `"tls-rf"` 或 `"ttl-d"`
- 增加 `num_records` 以获得更多分片
- 尝试不同的 `send_interval` 值

### 性能问题

- 减少 `num_records`（尝试5-10而不是20+）
- 减少 `send_interval`（尝试"50ms"或"100ms"）
- 对信任的域名使用 `"direct"` 模式

## 性能调优

### 低延迟配置
```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 5,
    "num_segs": 2,
    "send_interval": "50ms"
  }
}
```

### 高安全性配置
```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 20,
    "num_segs": 5,
    "send_interval": "300ms"
  }
}
```

### 平衡配置
```json
{
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  }
}
```

## 完整工作流程

1. **创建配置文件**
   - 使用提供的示例作为模板
   - 根据需要调整策略
   
2. **集成到应用**
   - 在VPN服务中添加lumine启动/停止逻辑
   - 修改sing-box配置生成器添加lumine出站
   
3. **测试连接**
   - 启动VPN
   - 验证流量通过lumine路由
   - 检查是否能正常访问网站
   
4. **优化性能**
   - 根据网络情况调整分片参数
   - 为不同域名设置不同策略
   - 监控延迟和稳定性

## 参考文档

- [LUMINE_CHAIN_PROXY_GUIDE.md](LUMINE_CHAIN_PROXY_GUIDE.md) - 英文详细指南
- [LUMINE_INTEGRATION.md](LUMINE_INTEGRATION.md) - 集成文档
- [singbox_lumine_chain_example.json](singbox_lumine_chain_example.json) - 完整配置示例
- [lumine_config_advanced.json](lumine_config_advanced.json) - 高级lumine配置

## 技术支持

如有问题，请查看：
- 原始lumine项目: https://github.com/UcnacDx2/lumine
- NekoBoxForAndroid项目: https://github.com/MatsuriDayo/NekoBoxForAndroid
