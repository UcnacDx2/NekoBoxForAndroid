# Lumine 配置示例 / Lumine Configuration Examples

这个目录包含可直接使用的 Lumine 配置文件和 NekoBox 导入链接。

This directory contains ready-to-use Lumine configuration files and NekoBox import links.

## 📋 快速使用 / Quick Usage

### 1. 选择配置 / Choose a Configuration

| 配置文件 / Config File | 导入链接 / Import Link | 用途 / Purpose |
|------------------------|----------------------|----------------|
| `lumine_basic.json` | `import_basic.txt` | 基础配置，适合大多数情况 / Basic config for most cases |
| `lumine_google.json` | `import_google.txt` | 针对 Google/YouTube 优化 / Optimized for Google/YouTube |
| `lumine_strong.json` | `import_strong.txt` | 强混淆，适合严格审查环境 / Strong obfuscation for strict censorship |
| `lumine_fast.json` | `import_fast.txt` | 快速模式，最小延迟 / Fast mode with minimal latency |

### 2. 启动 Lumine / Start Lumine

#### 使用 adb (需要电脑) / Using adb (requires computer):

```bash
# 1. 推送配置文件到手机
adb push lumine_basic.json /sdcard/lumine_config.json

# 2. 启动 Lumine
adb shell /data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &
```

#### 使用 Termux (推荐) / Using Termux (recommended):

```bash
# 1. 安装 Termux (从 F-Droid)
# 2. 将配置文件复制到 /sdcard/lumine_config.json
# 3. 在 Termux 中运行:
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &
```

### 3. 在 NekoBox 中导入 / Import in NekoBox

**方法 1: 从剪切板导入 (推荐) / From Clipboard (Recommended)**

1. 打开对应的 `import_*.txt` 文件 / Open corresponding `import_*.txt` file
2. 复制内容到剪切板 / Copy content to clipboard
3. 在 NekoBox 中点击"从剪切板导入" / Click "Import from clipboard" in NekoBox

**方法 2: 手动添加 / Manual Add**

1. 在 NekoBox 中点击 "+" / Click "+" in NekoBox
2. 选择 "SOCKS5" / Select "SOCKS5"
3. 填写:
   - 服务器 / Server: `127.0.0.1`
   - 端口 / Port: `1080`
   - 名称 / Name: `Lumine`

### 4. 连接 / Connect

在 NekoBox 中选择导入的配置并点击连接。

Select the imported config in NekoBox and click connect.

## 📝 配置说明 / Configuration Details

### lumine_basic.json

- **用途 / Purpose**: 日常使用的基础配置 / Basic config for daily use
- **特点 / Features**:
  - TLS 分片: 1024 字节 / TLS fragmentation: 1024 bytes
  - 分片延迟: 10ms / Fragment delay: 10ms
  - DNS over HTTPS
  - 适中的混淆效果 / Moderate obfuscation

### lumine_google.json

- **用途 / Purpose**: 专门针对 Google 服务优化 / Optimized for Google services
- **特点 / Features**:
  - Google/YouTube 域名特殊处理 / Special handling for Google/YouTube domains
  - 更大的分片尺寸提升速度 / Larger fragment size for better speed
  - 更低的延迟 / Lower latency
  - Google DNS (8.8.8.8)

### lumine_strong.json

- **用途 / Purpose**: 应对严格的网络审查 / For strict network censorship
- **特点 / Features**:
  - 小分片尺寸 (512 字节) / Small fragment size (512 bytes)
  - 更高的延迟以增强混淆 / Higher latency for better obfuscation
  - TCP 脱同步 / TCP desynchronization
  - 适合突破高级 DPI / Good for bypassing advanced DPI

### lumine_fast.json

- **用途 / Purpose**: 追求速度和低延迟 / For speed and low latency
- **特点 / Features**:
  - 大分片尺寸 (8192 字节) / Large fragment size (8192 bytes)
  - 最小延迟 (1ms) / Minimal delay (1ms)
  - 不进行分片处理 / No fragmentation processing
  - 同时提供 HTTP 代理 / Also provides HTTP proxy
  - 适合稳定网络环境 / Good for stable network environments

## 🔧 自定义配置 / Custom Configuration

可以基于这些配置文件进行自定义修改:

You can customize these configuration files:

### 修改端口 / Change Port

```json
{
  "socks5_address": "127.0.0.1:YOUR_PORT",
  ...
}
```

### 添加域名策略 / Add Domain Policies

```json
{
  ...
  "domain_policies": {
    "example.com;*.example.com": {
      "mode": "proxy",
      "fragment_size": 2048,
      "fragment_sleep": "5ms"
    }
  }
}
```

### 修改 DNS 服务器 / Change DNS Server

```json
{
  ...
  "dns_addr": "https://8.8.8.8/dns-query",  // Google
  // or
  "dns_addr": "https://1.1.1.1/dns-query",  // Cloudflare
  // or
  "dns_addr": "8.8.8.8:53",  // UDP DNS
  ...
}
```

## ⚠️ 重要提示 / Important Notes

1. **必须先启动 Lumine** 才能在 NekoBox 中连接
   **Lumine must be started first** before connecting in NekoBox

2. **检查端口是否可用** (默认 1080)
   **Check if port is available** (default 1080)
   ```bash
   netstat -an | grep 1080
   ```

3. **确保配置文件是有效的 JSON 格式**
   **Ensure config file is valid JSON format**

4. **某些设备可能需要 root 权限** 来访问 NekoBox 的 lib 目录
   **Some devices may require root access** to access NekoBox lib directory

## 🔗 相关文档 / Related Documentation

- [LUMINE_QUICK_START.md](../../LUMINE_QUICK_START.md) - 完整快速开始指南 / Complete quick start guide
- [LUMINE_USER_GUIDE.md](../../LUMINE_USER_GUIDE.md) - 详细用户手册 / Detailed user manual
- [LUMINE_READY_TO_USE_CONFIGS.txt](../../LUMINE_READY_TO_USE_CONFIGS.txt) - 纯文本格式配置 / Plain text configs

## 📞 获取帮助 / Get Help

- GitHub Issues: https://github.com/UcnacDx2/NekoBoxForAndroid
- Lumine Project: https://github.com/UcnacDx2/lumine
