# Lumine 网络预处理插件使用文档 / Lumine Network Preprocessing Plugin User Guide

[中文](#中文文档) | [English](#english-documentation)

---

## 中文文档

### 简介

Lumine 是一个轻量级的本地 HTTP/SOCKS5 代理服务器，基于 [TlsFragment](https://github.com/maoist2009/TlsFragment) 实现，专门用于保护 TCP 上的 TLS 连接，绕过网络审查。

### 主要功能

- **SOCKS5 代理服务器**：提供标准 SOCKS5 代理功能
- **HTTP 代理服务器**：支持 HTTP/HTTPS 代理
- **TLS 流量处理**：
  - TLS 记录分片 (TLS Record Fragmentation)
  - TCP 分片 (TCP Segmentation)
  - TTL 脱同步 (TTL Desynchronization)
  - 伪造 DNS 响应
- **灵活的策略配置**：支持基于域名、IP 的不同处理策略

### 安装方式

Lumine 已作为内置插件集成到 NekoBox for Android 中，无需额外下载。

### 构建方法

如果您需要自行编译：

```bash
cd libcore
./build_lumine.sh
```

这将为以下 Android 架构生成 `liblumine.so` 文件：
- ARM64 (arm64-v8a)
- ARMv7 (armeabi-v7a)
- x86
- x86_64

编译产物将输出到 `app/src/main/jniLibs/{arch}/liblumine.so`

### 配置说明

#### 基础配置

创建配置文件 `config.json`：

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "127.0.0.1:1081",
    "dns_addr": "https://1.1.1.1/dns-query",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3
    }
}
```

#### 配置字段说明

##### 顶层字段

| 字段 | 说明 | 示例 | 特殊值 |
|------|------|------|--------|
| `socks5_address` | SOCKS5 绑定地址 | `"127.0.0.1:1080"` | `"none"` 禁用 SOCKS5 |
| `http_address` | HTTP 绑定地址 | `":1081"` | `"none"` 禁用 HTTP |
| `dns_addr` | DNS 服务器地址 (UDP/HTTPS) | `"127.0.0.1:8053"`, `"https://1.1.1.1/dns-query"` | - |
| `udp_minsize` | DNS 查询的最小 UDP 包大小 | `4096` | `0` 使用默认大小 |
| `socks5_for_doh` | DoH 使用的 SOCKS5 代理 | `"127.0.0.1:1080"` | 空字符串禁用 |
| `max_jump` | IP 映射的最大重定向链长度 | `30` | `0` 默认为 20 |
| `fake_ttl_rules` | 伪包的 TTL 计算规则 | `"0-1;3=3;5-1;8-2"` | 空字符串禁用 |
| `transmit_file_limit` | 最大并发 TransmitFile 操作数 | `2` | `0` 或负数表示无限制 |
| `dns_cache_ttl` | DNS 答案缓存时长 (秒) | `259200` | `-1` 永久缓存；`0` 禁用缓存 |
| `fake_ttl_cache_ttl` | 最小可达 TTL 缓存时长 (秒) | `259200` | `-1` 永久缓存；`0` 禁用缓存 |

##### 策略字段

| 字段 | 说明 | 示例 | 特殊值 |
|------|------|------|--------|
| `connect_timeout` | 连接超时时间 | `"10s"` | - |
| `reply_first` | 在连接建立前发送 SOCKS5 成功响应 | `true` | - |
| `host` | 覆盖目标主机 | `"^208.103.161.2"`, `"www.ietf.org"` | 前缀 `^` 禁用 IP 重定向 |
| `map_to` | IP 重定向到其他主机/CIDR | `"35.180.16.12"`, `"^www.fbi.org"` | 前缀 `^` 禁用链式跳转 |
| `port` | 覆盖目标端口 | `8443` | `0` 使用原始端口 |
| `mode` | 流量处理模式 | `"tls-rf"` | 见下表 |
| `num_records` | TLS 分片的记录数 | `10` | `1` 禁用分片 |
| `num_segs` | TCP 分片的段数 | `3` | `1` 禁用分段；`-1` 每次发送 1 条记录 |
| `oob` | 在第一个 TCP 段末尾附加 OOB 数据 | `true` | - |
| `send_interval` | 段间发送间隔 | `"200ms"` | `0s` 无延迟 |
| `fake_ttl` | `ttl-d` 模式下伪包的 TTL 值 | `17` | `0` 启用自动 TTL 检测 |
| `fake_sleep` | 发送伪包后的休眠时间 | `"200ms"` | - |

##### 流量处理模式

| 模式 | 说明 | 用途 |
|------|------|------|
| `raw` | SOCKS5 后的原始 TCP 转发 | 最小开销 |
| `direct` | 不做处理的直通 | 一般流量 |
| `tls-rf` | TLS 记录分片 | TLS 连接 |
| `ttl-d` | 基于 TTL 的脱同步与伪包 | TLS 连接 |
| `block` | 完全阻止连接 | 连接终止 |
| `tls-alert` | 发送 TLS 警报并终止连接 | TLS 连接终止 |

### 使用示例

#### 示例 1：基础 TLS 分片配置

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "127.0.0.1:1081",
    "dns_addr": "https://1.1.1.1/dns-query",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3,
        "send_interval": "100ms"
    }
}
```

#### 示例 2：针对特定域名的 TTL 脱同步

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "none",
    "dns_addr": "https://cloudflare-dns.com/dns-query",
    "default_policy": {
        "mode": "direct"
    },
    "domain_policies": {
        "example.com": {
            "mode": "ttl-d",
            "fake_ttl": 0,
            "fake_sleep": "200ms"
        },
        "*.blocked-site.com": {
            "mode": "tls-rf",
            "num_records": 15,
            "num_segs": 5
        }
    }
}
```

#### 示例 3：IP 重定向和阻止

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "127.0.0.1:1081",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10
    },
    "ip_policies": {
        "8.8.8.8": {
            "map_to": "1.1.1.1"
        },
        "192.168.1.0/24": {
            "mode": "block"
        }
    }
}
```

### 在 NekoBox 中使用

1. **启动 Lumine**：
   - Lumine 作为内置插件，无需额外安装
   - 在 NekoBox 设置中找到 "Lumine" 选项
   - 配置 Lumine 的配置文件路径

2. **配置文件位置**：
   - 推荐放在：`/data/data/io.nekohasekai.sagernet/files/lumine_config.json`
   - 或使用外部存储：`/sdcard/NekoBox/lumine_config.json`

3. **与 NekoBox 集成**：
   - Lumine 可以作为前置代理使用
   - 将 NekoBox 的上游代理设置指向 Lumine (127.0.0.1:1080)
   - 或在路由规则中选择性地使用 Lumine

### 故障排除

#### 问题：Lumine 无法启动

- 检查配置文件是否有效的 JSON 格式
- 验证端口是否已被占用
- 查看日志文件获取详细错误信息

#### 问题：连接超时

- 增加 `connect_timeout` 值
- 检查 DNS 配置是否正确
- 尝试不同的流量处理模式

#### 问题：某些网站无法访问

- 尝试调整 `num_records` 和 `num_segs` 参数
- 使用 `ttl-d` 模式替代 `tls-rf`
- 检查域名策略配置

### 技术原理

Lumine 通过以下技术绕过网络审查：

1. **TLS 记录分片**：将 TLS 握手分成多个小记录，避免 DPI 检测
2. **TCP 分片**：在 TCP 层面进一步分片，增加检测难度
3. **TTL 脱同步**：发送低 TTL 的伪包欺骗中间设备
4. **DNS 操作**：通过 DoH 避免 DNS 污染

### 性能建议

- **num_records**: 5-15 之间通常效果最好
- **num_segs**: 3-5 可以在性能和效果之间平衡
- **send_interval**: 50-200ms，过大会影响速度
- **dns_cache_ttl**: 建议至少 3600 (1小时) 以减少 DNS 查询

---

## English Documentation

### Introduction

Lumine is a lightweight local HTTP/SOCKS5 proxy server based on [TlsFragment](https://github.com/maoist2009/TlsFragment) that protects TLS connections over TCP to bypass censorship.

### Key Features

- **SOCKS5 Proxy Server**: Standard SOCKS5 proxy functionality
- **HTTP Proxy Server**: HTTP/HTTPS proxy support
- **TLS Traffic Manipulation**:
  - TLS Record Fragmentation
  - TCP Segmentation
  - TTL Desynchronization
  - Fake DNS Responses
- **Flexible Policy Configuration**: Different handling strategies based on domains and IPs

### Installation

Lumine is integrated as a built-in plugin in NekoBox for Android, no additional download required.

### Building

To compile manually:

```bash
cd libcore
./build_lumine.sh
```

This generates `liblumine.so` files for Android architectures:
- ARM64 (arm64-v8a)
- ARMv7 (armeabi-v7a)
- x86
- x86_64

Output location: `app/src/main/jniLibs/{arch}/liblumine.so`

### Configuration

#### Basic Configuration

Create `config.json`:

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "127.0.0.1:1081",
    "dns_addr": "https://1.1.1.1/dns-query",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3
    }
}
```

#### Configuration Fields

##### Top-Level Fields

| Field | Description | Example | Special Values |
|-------|-------------|---------|----------------|
| `socks5_address` | SOCKS5 bind address | `"127.0.0.1:1080"` | `"none"` disables SOCKS5 |
| `http_address` | HTTP bind address | `":1081"` | `"none"` disables HTTP |
| `dns_addr` | DNS server address (UDP/HTTPS) | `"127.0.0.1:8053"`, `"https://1.1.1.1/dns-query"` | - |
| `udp_minsize` | Minimum UDP packet size for DNS | `4096` | `0` uses default size |
| `socks5_for_doh` | SOCKS5 proxy for DoH | `"127.0.0.1:1080"` | Empty string disables |
| `max_jump` | Maximum redirect chain length | `30` | `0` defaults to 20 |
| `fake_ttl_rules` | TTL rules for fake packets | `"0-1;3=3;5-1;8-2"` | Empty string disables |
| `transmit_file_limit` | Max concurrent TransmitFile ops | `2` | `0` or negative means unlimited |
| `dns_cache_ttl` | DNS cache duration (seconds) | `259200` | `-1` forever; `0` disabled |
| `fake_ttl_cache_ttl` | TTL cache duration (seconds) | `259200` | `-1` forever; `0` disabled |

##### Policy Fields

| Field | Description | Example | Special Values |
|-------|-------------|---------|----------------|
| `connect_timeout` | Connection timeout | `"10s"` | - |
| `reply_first` | Send SOCKS5 success before connecting | `true` | - |
| `host` | Override target host | `"^208.103.161.2"`, `"www.ietf.org"` | Prefix `^` disables IP redirection |
| `map_to` | Redirect IP to another host | `"35.180.16.12"`, `"^www.fbi.org"` | Prefix `^` disables chain jump |
| `port` | Override target port | `8443` | `0` uses original port |
| `mode` | Traffic manipulation mode | `"tls-rf"` | See table below |
| `num_records` | Number of TLS records | `10` | `1` disables fragmentation |
| `num_segs` | Number of TCP segments | `3` | `1` disables splitting; `-1` sends 1 record each time |
| `oob` | Attach OOB data to first segment | `true` | - |
| `send_interval` | Interval between segments | `"200ms"` | `0s` means no delay |
| `fake_ttl` | TTL for fake packets in `ttl-d` mode | `17` | `0` enables auto detection |
| `fake_sleep` | Sleep time after fake packet | `"200ms"` | - |

##### Traffic Modes

| Mode | Description | Use Case |
|------|-------------|----------|
| `raw` | Raw TCP forwarding | Minimal overhead |
| `direct` | Pass-through | General traffic |
| `tls-rf` | TLS record fragmentation | TLS connections |
| `ttl-d` | TTL desynchronization | TLS connections |
| `block` | Block connection | Connection termination |
| `tls-alert` | Send TLS alert | TLS termination |

### Usage Examples

#### Example 1: Basic TLS Fragmentation

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "127.0.0.1:1081",
    "dns_addr": "https://1.1.1.1/dns-query",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10,
        "num_segs": 3,
        "send_interval": "100ms"
    }
}
```

#### Example 2: Domain-Specific TTL Desync

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "none",
    "dns_addr": "https://cloudflare-dns.com/dns-query",
    "default_policy": {
        "mode": "direct"
    },
    "domain_policies": {
        "example.com": {
            "mode": "ttl-d",
            "fake_ttl": 0,
            "fake_sleep": "200ms"
        },
        "*.blocked-site.com": {
            "mode": "tls-rf",
            "num_records": 15,
            "num_segs": 5
        }
    }
}
```

#### Example 3: IP Redirection and Blocking

```json
{
    "socks5_address": "127.0.0.1:1080",
    "http_address": "127.0.0.1:1081",
    "default_policy": {
        "mode": "tls-rf",
        "num_records": 10
    },
    "ip_policies": {
        "8.8.8.8": {
            "map_to": "1.1.1.1"
        },
        "192.168.1.0/24": {
            "mode": "block"
        }
    }
}
```

### Using with NekoBox

1. **Start Lumine**:
   - Lumine is built-in, no extra installation needed
   - Find "Lumine" option in NekoBox settings
   - Configure the config file path

2. **Config File Location**:
   - Recommended: `/data/data/io.nekohasekai.sagernet/files/lumine_config.json`
   - Or external storage: `/sdcard/NekoBox/lumine_config.json`

3. **Integration with NekoBox**:
   - Use Lumine as an upstream proxy
   - Point NekoBox upstream to Lumine (127.0.0.1:1080)
   - Or selectively use Lumine in routing rules

### Troubleshooting

#### Issue: Lumine won't start

- Check if config file is valid JSON
- Verify ports are not in use
- Check logs for detailed errors

#### Issue: Connection timeouts

- Increase `connect_timeout` value
- Verify DNS configuration
- Try different traffic modes

#### Issue: Some sites don't work

- Adjust `num_records` and `num_segs` parameters
- Try `ttl-d` mode instead of `tls-rf`
- Check domain policy configuration

### Technical Principles

Lumine bypasses censorship using:

1. **TLS Record Fragmentation**: Split TLS handshake into small records to avoid DPI
2. **TCP Segmentation**: Further split at TCP layer to increase detection difficulty
3. **TTL Desynchronization**: Send low-TTL fake packets to deceive middleboxes
4. **DNS Manipulation**: Use DoH to avoid DNS poisoning

### Performance Tips

- **num_records**: 5-15 usually works best
- **num_segs**: 3-5 balances performance and effectiveness
- **send_interval**: 50-200ms, too large affects speed
- **dns_cache_ttl**: At least 3600 (1 hour) recommended to reduce DNS queries

---

## License

GPL-3.0

## Credits

- Original lumine implementation: https://github.com/UcnacDx2/lumine
- Based on TlsFragment: https://github.com/maoist2009/TlsFragment
- Integration by: NekoBox for Android contributors
