# Lumine VPN 快速配置指南 / Lumine VPN Quick Start Guide

[中文](#中文) | [English](#english)

---

## 中文

### 📋 可从剪切板导入的 Lumine SOCKS5 配置

#### 步骤 1: 启动 Lumine 服务

首先需要创建 Lumine 配置文件并启动服务。

**创建配置文件** `/sdcard/lumine_config.json`:

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "udp_minsize": 0,
  "dns_cache_ttl": 3600,
  "ttl_cache_ttl": 0,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "read_timeout": "30s",
    "write_timeout": "30s",
    "reply_first": false,
    "fragment_size": 1024,
    "fragment_sleep": "10ms",
    "fragment_type": "tls",
    "desync_zero": false,
    "desync_split": 2,
    "desync_split_position": 3,
    "desync_ttl": 0,
    "desync_disorder": false,
    "desync_fake_ttl": 0,
    "desync_fake_method": "",
    "desync_fake_data": "",
    "port": 0
  },
  "domain_policies": {},
  "ip_policies": {}
}
```

**启动 Lumine:**
```bash
# 方法 1: 使用 adb (需要电脑)
adb shell /data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &

# 方法 2: 使用终端应用 (Termux 等)
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &
```

#### 步骤 2: 从剪切板导入配置

**复制以下 SOCKS5 链接到剪切板:**

```
socks5://127.0.0.1:1080#Lumine-Local
```

**然后在 NekoBox 中:**
1. 打开 NekoBox 应用
2. 点击 "+" 或从剪切板导入按钮
3. 自动识别并导入 Lumine SOCKS5 配置
4. 点击连接

### 🔧 带用户名密码的 SOCKS5 配置

如果需要认证（建议本地使用时不需要）:

```
socks5://username:password@127.0.0.1:1080#Lumine-Auth
```

### 🚀 高级配置示例

#### 1. 针对 Google/YouTube 优化的配置

Lumine 配置文件 (`/sdcard/lumine_google.json`):
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://8.8.8.8/dns-query",
  "dns_cache_ttl": 3600,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "fragment_size": 2048,
    "fragment_sleep": "5ms",
    "fragment_type": "tls"
  },
  "domain_policies": {
    "google.com;*.google.com;*.googleapis.com": {
      "mode": "proxy",
      "fragment_size": 3072,
      "fragment_sleep": "3ms"
    },
    "youtube.com;*.youtube.com;*.googlevideo.com": {
      "mode": "proxy",
      "fragment_size": 4096,
      "fragment_sleep": "2ms"
    }
  }
}
```

**导入链接:**
```
socks5://127.0.0.1:1080#Lumine-Google
```

#### 2. 强混淆配置（适用于严格审查）

Lumine 配置文件 (`/sdcard/lumine_strong.json`):
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "dns_cache_ttl": 3600,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "fragment_size": 512,
    "fragment_sleep": "20ms",
    "fragment_type": "tls",
    "desync_split": 3,
    "desync_split_position": 2
  }
}
```

**导入链接:**
```
socks5://127.0.0.1:1080#Lumine-Strong
```

#### 3. 快速模式（最小延迟）

Lumine 配置文件 (`/sdcard/lumine_fast.json`):
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "dns_cache_ttl": 7200,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "5s",
    "fragment_size": 8192,
    "fragment_sleep": "1ms",
    "fragment_type": "none"
  }
}
```

**导入链接:**
```
socks5://127.0.0.1:1080#Lumine-Fast
```

### 📱 完整使用流程

#### 使用 Termux (推荐)

1. **安装 Termux** (从 F-Droid 下载)

2. **创建启动脚本** `/sdcard/start_lumine.sh`:
```bash
#!/system/bin/sh
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json
```

3. **在 Termux 中运行:**
```bash
sh /sdcard/start_lumine.sh &
```

4. **复制配置到剪切板:**
```
socks5://127.0.0.1:1080#Lumine
```

5. **在 NekoBox 中导入并连接**

### 🔄 链式代理配置

如果需要将 Lumine 与远程服务器链接使用:

**步骤:**
1. 先启动 Lumine (监听 127.0.0.1:1080)
2. 在 NekoBox 中添加两个配置:
   - 配置 1: `socks5://127.0.0.1:1080#Lumine-Local`
   - 配置 2: 你的远程服务器配置 (如 Shadowsocks, VMess 等)
3. 创建代理链: Lumine-Local → 远程服务器

### ⚠️ 注意事项

1. **Lumine 必须先启动** - 在连接 SOCKS5 配置前确保 Lumine 服务正在运行
2. **端口不要冲突** - 确保 1080 端口没有被其他程序占用
3. **权限问题** - 某些设备可能需要 root 权限才能运行 Lumine
4. **持久运行** - 使用 Termux 或其他方式让 Lumine 在后台持续运行

### 🛠️ 故障排除

**问题: 无法连接**
- 检查 Lumine 是否正在运行: `ps | grep lumine`
- 检查端口: `netstat -an | grep 1080`
- 查看 Lumine 日志

**问题: 速度慢**
- 减小 `fragment_size`
- 减小 `fragment_sleep`
- 尝试不同的 DNS 服务器

**问题: 某些网站无法访问**
- 调整 `fragment_type` (尝试 "tls", "tcp", "none")
- 增加 `desync_split` 参数
- 在 `domain_policies` 中为特定网站设置策略

---

## English

### 📋 Clipboard-Importable Lumine SOCKS5 Configuration

#### Step 1: Start Lumine Service

First, create Lumine config file and start the service.

**Create config file** `/sdcard/lumine_config.json`:

```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "udp_minsize": 0,
  "dns_cache_ttl": 3600,
  "ttl_cache_ttl": 0,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "read_timeout": "30s",
    "write_timeout": "30s",
    "reply_first": false,
    "fragment_size": 1024,
    "fragment_sleep": "10ms",
    "fragment_type": "tls",
    "desync_zero": false,
    "desync_split": 2,
    "desync_split_position": 3,
    "desync_ttl": 0,
    "desync_disorder": false,
    "desync_fake_ttl": 0,
    "desync_fake_method": "",
    "desync_fake_data": "",
    "port": 0
  },
  "domain_policies": {},
  "ip_policies": {}
}
```

**Start Lumine:**
```bash
# Method 1: Using adb (requires computer)
adb shell /data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &

# Method 2: Using terminal app (Termux, etc.)
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &
```

#### Step 2: Import from Clipboard

**Copy this SOCKS5 link to clipboard:**

```
socks5://127.0.0.1:1080#Lumine-Local
```

**Then in NekoBox:**
1. Open NekoBox app
2. Click "+" or import from clipboard button
3. Automatically recognize and import Lumine SOCKS5 config
4. Click connect

### 🔧 SOCKS5 with Authentication

If authentication is needed (not recommended for local use):

```
socks5://username:password@127.0.0.1:1080#Lumine-Auth
```

### 🚀 Advanced Configuration Examples

#### 1. Optimized for Google/YouTube

Lumine config file (`/sdcard/lumine_google.json`):
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://8.8.8.8/dns-query",
  "dns_cache_ttl": 3600,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "fragment_size": 2048,
    "fragment_sleep": "5ms",
    "fragment_type": "tls"
  },
  "domain_policies": {
    "google.com;*.google.com;*.googleapis.com": {
      "mode": "proxy",
      "fragment_size": 3072,
      "fragment_sleep": "3ms"
    },
    "youtube.com;*.youtube.com;*.googlevideo.com": {
      "mode": "proxy",
      "fragment_size": 4096,
      "fragment_sleep": "2ms"
    }
  }
}
```

**Import link:**
```
socks5://127.0.0.1:1080#Lumine-Google
```

#### 2. Strong Obfuscation (for strict censorship)

Lumine config file (`/sdcard/lumine_strong.json`):
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "dns_cache_ttl": 3600,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "fragment_size": 512,
    "fragment_sleep": "20ms",
    "fragment_type": "tls",
    "desync_split": 3,
    "desync_split_position": 2
  }
}
```

**Import link:**
```
socks5://127.0.0.1:1080#Lumine-Strong
```

#### 3. Fast Mode (minimal latency)

Lumine config file (`/sdcard/lumine_fast.json`):
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "dns_cache_ttl": 7200,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "5s",
    "fragment_size": 8192,
    "fragment_sleep": "1ms",
    "fragment_type": "none"
  }
}
```

**Import link:**
```
socks5://127.0.0.1:1080#Lumine-Fast
```

### 📱 Complete Usage Workflow

#### Using Termux (Recommended)

1. **Install Termux** (from F-Droid)

2. **Create start script** `/sdcard/start_lumine.sh`:
```bash
#!/system/bin/sh
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json
```

3. **Run in Termux:**
```bash
sh /sdcard/start_lumine.sh &
```

4. **Copy config to clipboard:**
```
socks5://127.0.0.1:1080#Lumine
```

5. **Import and connect in NekoBox**

### 🔄 Chain Proxy Configuration

To use Lumine with a remote server:

**Steps:**
1. Start Lumine first (listening on 127.0.0.1:1080)
2. Add two configs in NekoBox:
   - Config 1: `socks5://127.0.0.1:1080#Lumine-Local`
   - Config 2: Your remote server config (Shadowsocks, VMess, etc.)
3. Create proxy chain: Lumine-Local → Remote Server

### ⚠️ Important Notes

1. **Lumine must be started first** - Ensure Lumine service is running before connecting
2. **Avoid port conflicts** - Make sure port 1080 is not used by other programs
3. **Permission issues** - Some devices may require root access to run Lumine
4. **Persistent running** - Use Termux or other methods to keep Lumine running in background

### 🛠️ Troubleshooting

**Issue: Cannot connect**
- Check if Lumine is running: `ps | grep lumine`
- Check port: `netstat -an | grep 1080`
- View Lumine logs

**Issue: Slow speed**
- Reduce `fragment_size`
- Reduce `fragment_sleep`
- Try different DNS servers

**Issue: Some websites inaccessible**
- Adjust `fragment_type` (try "tls", "tcp", "none")
- Increase `desync_split` parameter
- Set policies for specific websites in `domain_policies`

---

## 🎯 Quick Reference

### Most Common Import Links

**Basic Lumine:**
```
socks5://127.0.0.1:1080#Lumine
```

**With custom name:**
```
socks5://127.0.0.1:1080#My-Lumine-Proxy
```

**Different port (if you configured Lumine on 1088):**
```
socks5://127.0.0.1:1088#Lumine-1088
```

### Configuration File Locations

- Main config: `/sdcard/lumine_config.json`
- Alternative: `/data/data/io.nekohasekai.sagernet/files/lumine_config.json`
- NekoBox data: `/data/data/io.nekohasekai.sagernet/`

---

## 📞 Support

For issues or questions:
- GitHub Issues: https://github.com/UcnacDx2/lumine
- NekoBox Issues: https://github.com/UcnacDx2/NekoBoxForAndroid

## 📄 License

GPL-3.0
