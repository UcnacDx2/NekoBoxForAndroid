# Lumine VPN - 在 NekoBox 中一键启用 / One-Click Enable in NekoBox

## 🎯 简单说明 / Simple Explanation

Lumine 是一个**本地代理插件**，它需要在本地运行并提供 SOCKS5 代理服务。NekoBox 通过连接到这个本地 SOCKS5 代理来使用 Lumine 的流量混淆功能。

这就像：
- **Lumine** = 一个在你手机上运行的小程序，提供流量混淆
- **NekoBox** = 连接到这个小程序来使用它的功能

Lumine is a **local proxy plugin** that runs locally and provides SOCKS5 proxy service. NekoBox connects to this local SOCKS5 proxy to use Lumine's traffic obfuscation features.

It's like:
- **Lumine** = A small program running on your phone that provides traffic obfuscation
- **NekoBox** = Connects to this program to use its features

---

## ✅ 最简单的使用方法 / Simplest Method

### 方案 A: 使用 Termux 自动化脚本 / Using Termux Automation Script

这个方法可以让 Lumine 在后台自动运行，NekoBox 只需要导入配置即可。

**步骤：**

1. **安装 Termux** (从 F-Droid 下载)

2. **创建自动启动脚本** `/sdcard/start_lumine.sh`:
```bash
#!/system/bin/sh
# 自动启动 Lumine
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &
```

3. **创建 Lumine 配置** `/sdcard/lumine_config.json`:
```json
{
  "socks5_address": "127.0.0.1:1080",
  "http_address": "none",
  "dns_addr": "https://1.1.1.1/dns-query",
  "dns_cache_ttl": 3600,
  "default_policy": {
    "mode": "proxy",
    "connect_timeout": "10s",
    "fragment_size": 1024,
    "fragment_sleep": "10ms",
    "fragment_type": "tls"
  },
  "domain_policies": {},
  "ip_policies": {}
}
```

4. **在 Termux 中运行**:
```bash
sh /sdcard/start_lumine.sh
```

5. **在 NekoBox 中导入**:
   - 复制: `socks5://127.0.0.1:1080#Lumine`
   - 在 NekoBox 点击"从剪切板导入"
   - 连接

**之后每次使用：**
- 打开 Termux，运行 `sh /sdcard/start_lumine.sh`
- 打开 NekoBox，连接 Lumine 配置

---

### 方案 B: 使用 Tasker/Automate 完全自动化 / Full Automation with Tasker/Automate

如果你想更自动化，可以使用 Tasker 或 Automate 应用：

1. **安装 Tasker** 或 **Automate**

2. **创建任务**：
   - 触发器: NekoBox 启动时
   - 动作: 运行 shell 命令 `/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &`

3. **在 NekoBox 中保存 Lumine 配置**（只需配置一次）

**之后：** 只需打开 NekoBox，Lumine 会自动启动！

---

## 🔧 为什么不能直接在 NekoBox 里配置？

### 技术原因 / Technical Reasons:

1. **Lumine 是独立进程** / Lumine runs as independent process
   - 类似 Hysteria、Naive 等插件
   - 需要单独的配置文件
   - 独立运行更稳定、安全

2. **架构限制** / Architecture Limitation
   - NekoBox 的插件系统设计为连接外部程序
   - Lumine 提供 SOCKS5 服务，NekoBox 作为客户端连接
   
3. **配置复杂性** / Configuration Complexity
   - Lumine 有很多高级选项（DNS、分片、脱同步等）
   - 独立配置文件更灵活、更强大

---

## 💡 未来的改进方案 / Future Improvement Options

如果要实现"在 NekoBox 内一键启用"，需要：

### 选项 1: 添加 Lumine 配置界面（需要大量开发）
- 在 NekoBox 中添加 Lumine 配置页面
- 自动生成配置文件
- 自动启动/停止 Lumine 进程
- **工作量**: 需要修改多个文件，添加新的 Activity

### 选项 2: 简化启动脚本（推荐的折中方案）
- 创建一个桌面快捷方式
- 一键启动 Lumine + NekoBox
- **工作量**: 很小，只需创建启动脚本

### 选项 3: 集成到 NekoBox（最彻底但最复杂）
- 将 Lumine 完全集成到 NekoBox 代码中
- 需要重新设计插件架构
- **工作量**: 非常大

---

## 🚀 推荐方案：创建桌面快捷方式 / Recommended: Desktop Shortcut

最实用的方案是创建一个启动快捷方式：

### 使用 Termux:Widget

1. **安装 Termux:Widget** (F-Droid)

2. **创建脚本** `~/.shortcuts/start-lumine-vpn.sh`:
```bash
#!/data/data/com.termux/files/usr/bin/sh
# 启动 Lumine
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &

# 等待 Lumine 启动
sleep 2

# 启动 NekoBox
am start -n io.nekohasekai.sagernet/.ui.MainActivity
```

3. **添加到桌面小部件**

**结果**: 点击桌面图标，Lumine 和 NekoBox 自动启动！

---

## 📋 总结 / Summary

| 方案 | 便捷度 | 复杂度 | 推荐度 |
|------|--------|--------|--------|
| 方案 A: Termux 手动 | ⭐⭐⭐ | ⭐ | ⭐⭐⭐⭐ |
| 方案 B: Tasker 自动化 | ⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| 桌面快捷方式 | ⭐⭐⭐⭐⭐ | ⭐ | ⭐⭐⭐⭐⭐ |
| 完全集成到 App | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐ (太复杂) |

**我的建议**: 使用 **Termux:Widget + 桌面快捷方式**，一键启动一切！

---

## ❓ 常见问题 / FAQ

**Q: 为什么其他 VPN 不需要这样？**  
A: 因为 Lumine 是专门的流量混淆工具，不是标准VPN协议。它提供更高级的功能。

**Q: 能不能自动启动 Lumine？**  
A: 可以！使用 Tasker 或 Termux:Widget 的桌面快捷方式。

**Q: Lumine 会一直运行吗？**  
A: 是的，它会在后台运行。如果想停止，运行 `pkill -f liblumine.so`

**Q: 每次重启手机都要重新启动吗？**  
A: 是的，或者使用 Tasker 设置开机自动启动。

---

## 📞 需要帮助？

如果你想要完整的自动化方案，我可以帮你创建：
1. 完整的 Termux 启动脚本
2. Tasker 配置文件
3. 桌面快捷方式

只需告诉我你想要哪种方案！
