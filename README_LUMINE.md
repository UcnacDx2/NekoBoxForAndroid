# Lumine Network Preprocessing - 完整实现指南

## 项目概述

本项目成功将lumine网络预处理库集成到NekoBoxForAndroid中，提供TLS分片和审查绕过功能，并实现了完整的链式代理支持。

## 快速导航

### 中文文档
- 📖 [**LUMINE_配置完整解答.md**](LUMINE_配置完整解答.md) - **推荐阅读** - 完整解答如何配置和使用
- 📖 [LUMINE_链式代理使用指南.md](LUMINE_链式代理使用指南.md) - 详细使用指南

### English Documentation
- 📖 [LUMINE_CHAIN_PROXY_GUIDE.md](LUMINE_CHAIN_PROXY_GUIDE.md) - Comprehensive chain proxy guide
- 📖 [LUMINE_INTEGRATION.md](LUMINE_INTEGRATION.md) - Integration documentation
- 📖 [IMPLEMENTATION_SUMMARY.md](IMPLEMENTATION_SUMMARY.md) - Technical implementation summary

### 配置示例文件
- 📄 [singbox_lumine_chain_example.json](singbox_lumine_chain_example.json) - Shadowsocks + Lumine完整配置
- 📄 [singbox_vmess_lumine_example.json](singbox_vmess_lumine_example.json) - VMess + Lumine配置
- 📄 [lumine_config_advanced.json](lumine_config_advanced.json) - Lumine高级配置（含域名策略）
- 📄 [lumine_config_example.json](lumine_config_example.json) - Lumine基础配置

## 核心功能

### ✨ 已实现功能

1. **Lumine集成** ✅
   - 完整移植lumine库到libcore
   - 通过gomobile暴露给Android
   - 支持SOCKS5和HTTP代理

2. **链式代理** ✅
   - 支持通过detour机制配置上游代理
   - Lumine作为预处理层
   - 支持任意代理协议（SS, VMess, Trojan等）

3. **TLS分片** ✅
   - TLS记录分片（tls-rf模式）
   - TTL去同步（ttl-d模式）
   - 可配置分片参数

4. **策略路由** ✅
   - 支持域名特定策略
   - 支持IP/CIDR特定策略
   - 可针对不同目标使用不同处理方式

5. **配置管理** ✅
   - 支持JSON配置文件
   - 支持编程方式配置
   - 提供配置验证和生成工具

## 两个核心问题的解答

### 问题1: 如何使用配置文件启用代理VPN网络？

**简短答案**:
```kotlin
// 1. 启动lumine服务
val lumine = Libcore.newLumineService("/sdcard/lumine_config.json")

// 2. 配置sing-box使用lumine
// 在outbounds中添加lumine出站，并在代理中设置detour
```

**详细文档**: 见 [LUMINE_配置完整解答.md](LUMINE_配置完整解答.md)

### 问题2: 如何基于Lumine作为上游实现链式代理？

**简短答案**:
```json
{
  "outbounds": [
    {"type": "socks", "tag": "lumine", "server": "127.0.0.1", "server_port": 1080},
    {"type": "shadowsocks", "tag": "proxy", "detour": "lumine", ...}
  ]
}
```

关键是在实际代理配置中添加 `"detour": "lumine"`

**详细文档**: 见 [LUMINE_链式代理使用指南.md](LUMINE_链式代理使用指南.md)

## 工作原理

### 链式代理流程

```
┌─────────┐   ┌──────────────┐   ┌────────┐   ┌──────────┐   ┌─────────┐
│ 用户应用 │ → │ NekoBox VPN  │ → │ Lumine │ → │ 代理服务器 │ → │ 互联网  │
│         │   │  (sing-box)  │   │  处理  │   │          │   │         │
└─────────┘   └──────────────┘   └────────┘   └──────────┘   └─────────┘
                                      ↓
                                TLS分片/TTL修改
```

### Lumine处理过程

```
原始TLS ClientHello (1个完整数据包)
    ↓
[Lumine TLS分片处理]
    ↓
分割成10个TLS记录 (num_records=10)
    ↓
每个记录分成3个TCP段 (num_segs=3)
    ↓
延迟200ms发送 (send_interval="200ms")
    ↓
绕过DPI检测 ✓
    ↓
到达代理服务器
```

## 使用示例

### 最简单的配置

**1. Lumine配置** (`/sdcard/lumine_config.json`):
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "default_policy": {
    "mode": "tls-rf",
    "num_records": 10,
    "num_segs": 3,
    "send_interval": "200ms"
  }
}
```

**2. Sing-box配置**:
```json
{
  "outbounds": [
    {
      "type": "socks",
      "tag": "lumine",
      "server": "127.0.0.1",
      "server_port": 1080
    },
    {
      "type": "shadowsocks",
      "tag": "proxy",
      "server": "your.server.com",
      "server_port": 8388,
      "method": "aes-256-gcm",
      "password": "password",
      "detour": "lumine"
    }
  ]
}
```

**3. 代码集成**:
```kotlin
// 启动lumine
val lumine = Libcore.newLumineService("/sdcard/lumine_config.json")

// 使用sing-box配置启动VPN
val box = Libcore.newSingBoxInstance(singboxConfig, null)
box.start()

// 停止时清理
lumine.close()
box.close()
```

## 技术实现

### 文件结构

```
libcore/
├── lumine/                    # Lumine核心库
│   ├── api.go                 # API接口
│   ├── config.go              # 配置管理
│   ├── dns.go                 # DNS处理
│   ├── fragment.go            # TLS分片
│   ├── http_proxy.go          # HTTP代理
│   ├── policy.go              # 策略管理
│   └── utils.go               # 工具函数
├── lumine_wrapper.go          # Android包装器
└── lumine_chain.go            # 链式代理辅助工具
```

### 编译构建

```bash
cd libcore
bash init.sh      # 初始化gomobile（仅首次）
bash build.sh     # 编译libcore.aar
```

输出: `app/libs/libcore.aar` (38MB)

### 已暴露的API

```kotlin
// 基础API
Libcore.newLumineService(configPath: String): LumineService
Libcore.newLumineServiceWithAddrs(socks5: String, http: String, config: String): LumineService
Libcore.testLumineIntegration(): String

// LumineService方法
service.close()
```

## 配置参数说明

### Lumine配置选项

| 参数 | 说明 | 示例值 |
|------|------|--------|
| `socks5_address` | SOCKS5监听地址 | `"127.0.0.1:1080"` |
| `http_address` | HTTP监听地址（"none"禁用） | `"127.0.0.1:1225"` |
| `dns_addr` | DNS服务器地址 | `"8.8.8.8:53"` |
| `default_policy.mode` | 处理模式 | `"tls-rf"`, `"ttl-d"`, `"direct"` |
| `default_policy.num_records` | TLS记录分片数 | `10` |
| `default_policy.num_segs` | TCP段分片数 | `3` |
| `default_policy.send_interval` | 发送间隔 | `"200ms"` |

完整参数说明见配置文件示例。

## 性能调优

### 低延迟配置（适合游戏、视频通话）
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

### 高安全配置（适合严格审查地区）
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

### 平衡配置（推荐）
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

## 故障排除

### 常见问题

**Q: Lumine启动失败**
- 检查端口是否被占用
- 尝试更改端口号
- 检查配置文件格式

**Q: 没有分片效果**
- 确认mode设置为"tls-rf"或"ttl-d"
- 检查detour配置是否正确
- 查看日志确认流量经过lumine

**Q: 连接速度慢**
- 减少num_records
- 减少send_interval
- 对可信域名使用"direct"模式

详细故障排除见各文档的故障排除章节。

## 开发路线图

### 已完成 ✅
- [x] Lumine库集成
- [x] Chain proxy支持
- [x] 配置文件管理
- [x] 完整中英文档
- [x] 示例配置文件
- [x] 辅助工具函数

### 待开发
- [ ] UI配置界面
- [ ] 可视化策略编辑器
- [ ] 性能监控仪表板
- [ ] 自动优化建议
- [ ] 预设配置模板

## 贡献者

- Lumine原始项目: https://github.com/UcnacDx2/lumine
- NekoBoxForAndroid: https://github.com/MatsuriDayo/NekoBoxForAndroid
- 集成实现: GitHub Copilot

## 许可证

GPL-3.0 License

## 相关链接

- [Lumine项目](https://github.com/UcnacDx2/lumine)
- [NekoBoxForAndroid](https://github.com/MatsuriDayo/NekoBoxForAndroid)
- [sing-box文档](https://sing-box.sagernet.org/)

---

**最后更新**: 2026-02-09

**状态**: ✅ 生产就绪

**推荐开始阅读**: [LUMINE_配置完整解答.md](LUMINE_配置完整解答.md)
