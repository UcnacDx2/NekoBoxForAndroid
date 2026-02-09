# Lumine 增强代理链配置示例 / Lumine Enhanced Chain Proxy Example

[中文](#中文说明) | [English](#english-description)

---

## 中文说明

### 配置概述

这个示例展示如何使用 Lumine 插件作为第一层代理，结合远程代理服务器，实现流量混淆和深度包检测（DPI）绕过。

### 使用场景

- 绕过深度包检测（DPI）
- 避免 TLS 流量识别
- 增强连接稳定性
- 突破高级网络审查

### 代理流程

```
客户端
  ↓
Lumine 插件 (流量混淆) - 127.0.0.1:1080
  ↓
远程 Shadowsocks/VMess 服务器
  ↓
目标网站
```

### Lumine 功能特性

- **TLS 记录分片**: 将 TLS 握手分成多个小记录
- **TCP 分片**: 在 TCP 层面进一步分片
- **TTL 脱同步**: 发送低 TTL 的伪包欺骗中间设备
- **DNS over HTTPS**: 避免 DNS 污染

## 配置步骤

### 第一步：配置 Lumine

1. **创建 Lumine 配置文件** `/sdcard/NekoBox/lumine_config.json`:

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "none",
    "dns_addr": "https://1.1.1.1/dns-query",
    "udp_minsize": 4096,
    "dns_cache_ttl": 3600,
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3,
        "send_interval": "100ms",
        "connect_timeout": "10s"
    },
    "domain_policies": {
        "*.google.com": {
            "mode": "tls-rf",
            "num_records": 15,
            "num_segs": 5
        },
        "*.youtube.com": {
            "mode": "ttl-d",
            "fake_ttl": 0,
            "fake_sleep": "200ms"
        }
    }
}
```

2. **配置说明**:
   - `socks5_address`: Lumine 监听地址（本地）
   - `http_address`: HTTP 代理地址（设为 "none" 禁用）
   - `dns_addr`: DNS 服务器（使用 DoH）
   - `default_policy.mode`: 默认使用 TLS 记录分片模式
   - `num_records`: 分片数量（10-15 效果较好）
   - `num_segs`: TCP 分段数（3-5 平衡性能）

### 第二步：在 NekoBox 中添加 Lumine 代理

1. 打开 NekoBox
2. 点击 "+" 按钮
3. 选择 "SOCKS5"
4. 填写配置信息：
   - **名称**: Lumine-Local
   - **服务器地址**: 127.0.0.1
   - **端口**: 1080
   - **用户名**: (留空)
   - **密码**: (留空)
5. 点击保存

### 第三步：添加远程代理服务器

#### 选项 A：Shadowsocks

1. 点击 "+" 按钮
2. 选择 "Shadowsocks"
3. 填写配置信息：
   - **名称**: SS-Server
   - **服务器地址**: your-server.com
   - **端口**: 8388
   - **密码**: your-password
   - **加密方式**: aes-256-gcm
4. 点击保存

#### 选项 B：VMess

1. 点击 "+" 按钮
2. 选择 "VMess"
3. 填写配置信息：
   - **名称**: VMess-Server
   - **服务器地址**: your-server.com
   - **端口**: 443
   - **UUID**: your-uuid
   - **加密方式**: auto
   - **传输协议**: ws + tls
4. 点击保存

### 第四步：创建 Lumine 链式代理

1. 点击 "+" 按钮
2. 选择 "Proxy Chain" (链式代理)
3. 填写配置信息：
   - **名称**: Lumine-Enhanced-Chain
4. 点击 "+" 添加代理：
   - 第一个：选择 "Lumine-Local" (SOCKS5)
   - 第二个：选择 "SS-Server" 或 "VMess-Server"
5. 确认顺序：Lumine 在前，远程服务器在后
6. 点击保存

### 第五步：启动 Lumine 服务

在使用链式代理前，需要确保 Lumine 服务正在运行。

**方法 1：使用内置 Lumine**
1. 在 NekoBox 设置中找到 Lumine 插件
2. 指定配置文件路径
3. 启动 Lumine 服务

**方法 2：手动启动 Lumine**
```bash
# 如果你有终端访问权限
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/NekoBox/lumine_config.json
```

### 第六步：使用链式代理

1. 返回 NekoBox 主界面
2. 选择 "Lumine-Enhanced-Chain"
3. 点击连接按钮
4. 等待连接成功

## 高级配置

### 针对不同网站的策略

在 Lumine 配置中，可以为不同域名设置不同的流量处理策略：

```json
{
    "domain_policies": {
        "*.google.com": {
            "mode": "tls-rf",
            "num_records": 15,
            "num_segs": 5,
            "send_interval": "50ms"
        },
        "*.youtube.com": {
            "mode": "ttl-d",
            "fake_ttl": 0,
            "fake_sleep": "200ms"
        },
        "*.twitter.com": {
            "mode": "tls-rf",
            "num_records": 12,
            "num_segs": 4
        },
        "*.facebook.com": {
            "mode": "direct"
        }
    }
}
```

### IP 特定策略

也可以为特定 IP 范围设置策略：

```json
{
    "ip_policies": {
        "8.8.8.8": {
            "map_to": "1.1.1.1"
        },
        "192.168.0.0/16": {
            "mode": "direct"
        }
    }
}
```

## 流量模式说明

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| `tls-rf` | TLS 记录分片 | 大多数 HTTPS 网站 |
| `ttl-d` | TTL 脱同步 | 严格的 DPI 检测 |
| `direct` | 直接连接（不处理） | 不需要混淆的流量 |
| `raw` | 原始 TCP 转发 | 最小开销 |
| `block` | 阻止连接 | 屏蔽特定域名 |

## 性能调优

### 1. 分片参数优化

- **num_records**: 10-15（推荐）
  - 太小：效果不明显
  - 太大：性能下降

- **num_segs**: 3-5（推荐）
  - 3：平衡性能和效果
  - 5：更强的混淆，但延迟更高

- **send_interval**: 50-200ms
  - 50ms：速度优先
  - 200ms：稳定性优先

### 2. DNS 优化

```json
{
    "dns_addr": "https://1.1.1.1/dns-query",
    "dns_cache_ttl": 3600,
    "udp_minsize": 4096
}
```

- 使用 DoH 避免 DNS 污染
- 启用 DNS 缓存减少查询
- 设置合适的 UDP 包大小

### 3. 连接优化

```json
{
    "default_policy": {
        "connect_timeout": "10s",
        "reply_first": false
    }
}
```

## 故障排除

### 问题 1：Lumine 服务无法启动

**可能原因**:
- 配置文件格式错误
- 端口被占用
- 权限不足

**解决方法**:
1. 验证 JSON 格式是否正确
2. 检查端口 1080 是否可用
3. 查看 Lumine 日志

### 问题 2：连接速度很慢

**可能原因**:
- 分片参数设置过高
- send_interval 设置过大
- 服务器距离远

**解决方法**:
1. 减小 num_records 和 num_segs
2. 降低 send_interval 到 50-100ms
3. 选择地理位置更近的服务器

### 问题 3：某些网站无法访问

**可能原因**:
- 流量模式不适合该网站
- 域名策略配置错误

**解决方法**:
1. 尝试不同的 mode（tls-rf 或 ttl-d）
2. 调整 num_records 参数
3. 检查域名策略配置

### 问题 4：DNS 解析失败

**可能原因**:
- DoH 服务器不可用
- DNS 配置错误

**解决方法**:
1. 更换 DNS 服务器（例如：`https://8.8.8.8/dns-query`）
2. 检查 DNS 缓存设置
3. 尝试使用 UDP DNS

## 验证配置

### 1. 测试 Lumine 服务

```bash
# 使用 curl 测试 SOCKS5 代理
curl -x socks5://127.0.0.1:1080 https://www.google.com
```

### 2. 检查流量混淆

使用 Wireshark 或 tcpdump 抓包，观察：
- TLS 握手是否被分片
- TCP 数据包是否被分段
- 是否有伪造的 TTL 包

### 3. 验证 IP 地址

访问 IP 查询网站，确认：
- 显示的 IP 是远程服务器的 IP
- 不是你的真实 IP

## 最佳实践

### 1. 配置建议

- **一般网站**: 使用 `tls-rf` 模式，num_records=10
- **高度审查**: 使用 `ttl-d` 模式，启用自动 TTL 检测
- **视频流**: 使用 `tls-rf` 模式，减小 send_interval

### 2. 安全建议

- 定期更换 Lumine 参数
- 为不同网站使用不同策略
- 结合多层代理使用

### 3. 性能建议

- 监控延迟和速度
- 根据实际情况调整参数
- 使用 DNS 缓存减少查询

---

## English Description

### Configuration Overview

This example demonstrates how to use the Lumine plugin as the first layer proxy, combined with a remote proxy server, to achieve traffic obfuscation and Deep Packet Inspection (DPI) bypass.

### Use Cases

- Bypass Deep Packet Inspection (DPI)
- Avoid TLS traffic recognition
- Enhance connection stability
- Break through advanced network censorship

### Proxy Flow

```
Client
  ↓
Lumine Plugin (Traffic Obfuscation) - 127.0.0.1:1080
  ↓
Remote Shadowsocks/VMess Server
  ↓
Target Website
```

### Lumine Features

- **TLS Record Fragmentation**: Split TLS handshake into small records
- **TCP Segmentation**: Further split at TCP layer
- **TTL Desynchronization**: Send low-TTL fake packets to deceive middleboxes
- **DNS over HTTPS**: Avoid DNS poisoning

## Configuration Steps

### Step 1: Configure Lumine

1. **Create Lumine config file** `/sdcard/NekoBox/lumine_config.json`:

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "none",
    "dns_addr": "https://1.1.1.1/dns-query",
    "udp_minsize": 4096,
    "dns_cache_ttl": 3600,
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3,
        "send_interval": "100ms",
        "connect_timeout": "10s"
    },
    "domain_policies": {
        "*.google.com": {
            "mode": "tls-rf",
            "num_records": 15,
            "num_segs": 5
        },
        "*.youtube.com": {
            "mode": "ttl-d",
            "fake_ttl": 0,
            "fake_sleep": "200ms"
        }
    }
}
```

2. **Configuration Explanation**:
   - `socks5_address`: Lumine listen address (local)
   - `http_address`: HTTP proxy address (set "none" to disable)
   - `dns_addr`: DNS server (using DoH)
   - `default_policy.mode`: Default TLS record fragmentation mode
   - `num_records`: Number of fragments (10-15 works well)
   - `num_segs`: TCP segments (3-5 balances performance)

### Step 2: Add Lumine Proxy in NekoBox

1. Open NekoBox
2. Click "+" button
3. Select "SOCKS5"
4. Fill in configuration:
   - **Name**: Lumine-Local
   - **Server Address**: 127.0.0.1
   - **Port**: 1080
   - **Username**: (leave empty)
   - **Password**: (leave empty)
5. Click save

### Step 3: Add Remote Proxy Server

#### Option A: Shadowsocks

1. Click "+" button
2. Select "Shadowsocks"
3. Fill in configuration:
   - **Name**: SS-Server
   - **Server Address**: your-server.com
   - **Port**: 8388
   - **Password**: your-password
   - **Encryption**: aes-256-gcm
4. Click save

#### Option B: VMess

1. Click "+" button
2. Select "VMess"
3. Fill in configuration:
   - **Name**: VMess-Server
   - **Server Address**: your-server.com
   - **Port**: 443
   - **UUID**: your-uuid
   - **Encryption**: auto
   - **Transport**: ws + tls
4. Click save

### Step 4: Create Lumine Chain Proxy

1. Click "+" button
2. Select "Proxy Chain"
3. Fill in configuration:
   - **Name**: Lumine-Enhanced-Chain
4. Click "+" to add proxies:
   - First: Select "Lumine-Local" (SOCKS5)
   - Second: Select "SS-Server" or "VMess-Server"
5. Confirm order: Lumine first, remote server second
6. Click save

### Step 5: Start Lumine Service

Before using the chain proxy, ensure Lumine service is running.

**Method 1: Use Built-in Lumine**
1. Find Lumine plugin in NekoBox settings
2. Specify config file path
3. Start Lumine service

**Method 2: Manually Start Lumine**
```bash
# If you have terminal access
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/NekoBox/lumine_config.json
```

### Step 6: Use Chain Proxy

1. Return to NekoBox main screen
2. Select "Lumine-Enhanced-Chain"
3. Click connect button
4. Wait for connection success

## Advanced Configuration

### Policies for Different Websites

In Lumine config, you can set different traffic handling policies for different domains:

```json
{
    "domain_policies": {
        "*.google.com": {
            "mode": "tls-rf",
            "num_records": 15,
            "num_segs": 5,
            "send_interval": "50ms"
        },
        "*.youtube.com": {
            "mode": "ttl-d",
            "fake_ttl": 0,
            "fake_sleep": "200ms"
        },
        "*.twitter.com": {
            "mode": "tls-rf",
            "num_records": 12,
            "num_segs": 4
        },
        "*.facebook.com": {
            "mode": "direct"
        }
    }
}
```

### IP-Specific Policies

You can also set policies for specific IP ranges:

```json
{
    "ip_policies": {
        "8.8.8.8": {
            "map_to": "1.1.1.1"
        },
        "192.168.0.0/16": {
            "mode": "direct"
        }
    }
}
```

## Traffic Mode Description

| Mode | Description | Use Case |
|------|-------------|----------|
| `tls-rf` | TLS record fragmentation | Most HTTPS websites |
| `ttl-d` | TTL desynchronization | Strict DPI detection |
| `direct` | Direct connection (no processing) | Traffic that doesn't need obfuscation |
| `raw` | Raw TCP forwarding | Minimal overhead |
| `block` | Block connection | Block specific domains |

## Performance Tuning

### 1. Fragmentation Parameter Optimization

- **num_records**: 10-15 (recommended)
  - Too small: Not effective
  - Too large: Performance degradation

- **num_segs**: 3-5 (recommended)
  - 3: Balance performance and effectiveness
  - 5: Stronger obfuscation, but higher latency

- **send_interval**: 50-200ms
  - 50ms: Speed priority
  - 200ms: Stability priority

### 2. DNS Optimization

```json
{
    "dns_addr": "https://1.1.1.1/dns-query",
    "dns_cache_ttl": 3600,
    "udp_minsize": 4096
}
```

- Use DoH to avoid DNS poisoning
- Enable DNS caching to reduce queries
- Set appropriate UDP packet size

### 3. Connection Optimization

```json
{
    "default_policy": {
        "connect_timeout": "10s",
        "reply_first": false
    }
}
```

## Troubleshooting

### Issue 1: Lumine Service Won't Start

**Possible Causes**:
- Config file format error
- Port occupied
- Insufficient permissions

**Solutions**:
1. Verify JSON format is correct
2. Check if port 1080 is available
3. Check Lumine logs

### Issue 2: Very Slow Connection

**Possible Causes**:
- Fragmentation parameters set too high
- send_interval set too large
- Server far away

**Solutions**:
1. Reduce num_records and num_segs
2. Lower send_interval to 50-100ms
3. Choose geographically closer server

### Issue 3: Some Websites Inaccessible

**Possible Causes**:
- Traffic mode not suitable for the website
- Domain policy misconfigured

**Solutions**:
1. Try different modes (tls-rf or ttl-d)
2. Adjust num_records parameter
3. Check domain policy configuration

### Issue 4: DNS Resolution Failure

**Possible Causes**:
- DoH server unavailable
- DNS configuration error

**Solutions**:
1. Change DNS server (e.g., `https://8.8.8.8/dns-query`)
2. Check DNS cache settings
3. Try using UDP DNS

## Verify Configuration

### 1. Test Lumine Service

```bash
# Test SOCKS5 proxy using curl
curl -x socks5://127.0.0.1:1080 https://www.google.com
```

### 2. Check Traffic Obfuscation

Use Wireshark or tcpdump to capture packets and observe:
- TLS handshake fragmentation
- TCP packet segmentation
- Fake TTL packets

### 3. Verify IP Address

Visit IP lookup website to confirm:
- Displayed IP is remote server IP
- Not your real IP

## Best Practices

### 1. Configuration Recommendations

- **General websites**: Use `tls-rf` mode, num_records=10
- **Heavy censorship**: Use `ttl-d` mode, enable auto TTL detection
- **Video streaming**: Use `tls-rf` mode, reduce send_interval

### 2. Security Recommendations

- Regularly change Lumine parameters
- Use different policies for different websites
- Combine with multi-layer proxies

### 3. Performance Recommendations

- Monitor latency and speed
- Adjust parameters based on actual conditions
- Use DNS caching to reduce queries

---

## 许可证 / License

GPL-3.0
