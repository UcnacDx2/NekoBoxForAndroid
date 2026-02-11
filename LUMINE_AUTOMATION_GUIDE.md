# Lumine 自动化启动指南 / Lumine Automation Setup Guide

本指南提供多种方式实现 Lumine 的自动化启动，让你可以"一键连接"。

This guide provides multiple ways to automate Lumine startup for "one-click connection".

---

## 🎯 方案概览 / Solution Overview

| 方案 | 便捷度 | 难度 | 适用人群 |
|------|--------|------|----------|
| A. Termux手动启动 | ⭐⭐⭐ | ⭐ | 所有用户 |
| B. Termux Widget桌面快捷 | ⭐⭐⭐⭐⭐ | ⭐⭐ | 推荐！ |
| C. Tasker自动化 | ⭐⭐⭐⭐ | ⭐⭐⭐ | 高级用户 |
| D. Automate流程 | ⭐⭐⭐⭐ | ⭐⭐ | 可视化操作 |

---

## 方案 A: Termux 手动启动 (最简单)

### 1. 安装准备

**下载 Termux** (从 F-Droid):
```
https://f-droid.org/packages/com.termux/
```

### 2. 复制脚本到手机

将以下脚本保存为 `/sdcard/start_lumine_nekobox.sh`:

```bash
#!/system/bin/sh
CONFIG_PATH="/sdcard/lumine_config.json"

# 创建默认配置（如果不存在）
if [ ! -f "$CONFIG_PATH" ]; then
    cat > "$CONFIG_PATH" << 'LUMINECONFIG'
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
LUMINECONFIG
fi

# 启动 Lumine
/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c "$CONFIG_PATH" &
sleep 2

# 启动 NekoBox
am start -n io.nekohasekai.sagernet/.ui.MainActivity

echo "✓ Lumine + NekoBox 已启动"
echo "在 NekoBox 中导入: socks5://127.0.0.1:1080#Lumine"
```

### 3. 使用方法

每次使用时:
1. 打开 Termux
2. 运行: `sh /sdcard/start_lumine_nekobox.sh`
3. NekoBox 自动打开，连接 Lumine 配置

**优点**: 简单、可靠  
**缺点**: 需要手动打开 Termux

---

## 方案 B: Termux:Widget 桌面快捷方式 (强烈推荐!) ⭐⭐⭐⭐⭐

这个方案可以在桌面添加一个图标，点击即可启动！

### 1. 安装应用

从 F-Droid 安装:
- **Termux** (`com.termux`)
- **Termux:Widget** (`com.termux.widget`)

### 2. 创建快捷方式脚本

在 Termux 中运行:

```bash
# 创建脚本目录
mkdir -p ~/.shortcuts

# 创建启动脚本
cat > ~/.shortcuts/start-lumine-vpn.sh << 'EOF'
#!/data/data/com.termux/files/usr/bin/sh

CONFIG_PATH="/sdcard/lumine_config.json"

# 创建默认配置
if [ ! -f "$CONFIG_PATH" ]; then
    cat > "$CONFIG_PATH" << 'LUMINECONFIG'
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
LUMINECONFIG
fi

# 检查 Lumine 是否已运行
if ! pgrep -f "liblumine.so" > /dev/null; then
    /data/data/io.nekohasekai.sagernet/lib/liblumine.so -c "$CONFIG_PATH" &
    sleep 2
fi

# 启动 NekoBox
am start -n io.nekohasekai.sagernet/.ui.MainActivity

# 显示通知
termux-toast "✓ Lumine + NekoBox 已启动"
EOF

# 添加执行权限
chmod +x ~/.shortcuts/start-lumine-vpn.sh

# 创建停止脚本
cat > ~/.shortcuts/stop-lumine-vpn.sh << 'EOF'
#!/data/data/com.termux/files/usr/bin/sh

PID=$(pgrep -f "liblumine.so")
if [ ! -z "$PID" ]; then
    kill $PID
    termux-toast "✓ Lumine 已停止"
else
    termux-toast "Lumine 未运行"
fi
EOF

chmod +x ~/.shortcuts/stop-lumine-vpn.sh

echo "✓ 脚本创建完成！"
echo "现在可以添加桌面小部件了。"
```

### 3. 添加桌面小部件

1. 长按桌面空白处
2. 选择 "小部件" (Widgets)
3. 找到 "Termux:Widget"
4. 将小部件拖到桌面
5. 现在你会看到 "start-lumine-vpn" 按钮

### 4. 使用

**启动**: 点击桌面的 "start-lumine-vpn" 图标  
**停止**: 点击桌面的 "stop-lumine-vpn" 图标

**优点**: 
- ✓ 一键启动
- ✓ 桌面图标
- ✓ 自动检测已运行状态
- ✓ 有通知提示

**缺点**: 需要安装两个应用

---

## 方案 C: Tasker 自动化 (高级)

### 1. 安装 Tasker

从 Play Store 或 F-Droid 安装 Tasker

### 2. 创建任务

**任务 1: 启动 Lumine**

1. 打开 Tasker
2. Tasks 标签 → 点击 "+"
3. 命名为 "Start Lumine"
4. 添加动作:
   - Action: Run Shell
   - Command: `/data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json &`
   - Timeout: 10 seconds
   - Use Root: NO

**任务 2: 启动 NekoBox**

1. 创建新任务 "Start NekoBox"
2. 添加动作:
   - Action: Launch App
   - App: NekoBox

**任务 3: 组合任务**

1. 创建新任务 "Start Lumine VPN"
2. 添加动作:
   - Perform Task: Start Lumine
   - Wait: 2 seconds
   - Perform Task: Start NekoBox

### 3. 创建触发器 (可选)

**选项 A: 桌面图标**
1. 长按桌面 → Widgets → Tasker → Task Shortcut
2. 选择 "Start Lumine VPN"

**选项 B: 打开 NekoBox 时自动启动**
1. Profiles 标签 → 点击 "+"
2. Event → App → App Changed
3. Application: NekoBox
4. 返回，选择任务 "Start Lumine"

---

## 方案 D: Automate 流程 (可视化)

### 1. 安装 Automate

从 Play Store 下载: LlamaLab Automate

### 2. 创建流程

1. 打开 Automate，创建新流程
2. 添加模块:

```
START (Flow Began)
  ↓
SHELL COMMAND EXECUTE
  Command: /data/data/io.nekohasekai.sagernet/lib/liblumine.so -c /sdcard/lumine_config.json
  Background: YES
  ↓
DELAY
  Duration: 2 seconds
  ↓
APP START
  Package: io.nekohasekai.sagernet
  Activity: .ui.MainActivity
  ↓
TOAST SHOW
  Message: Lumine + NekoBox started
```

3. 保存流程为 "Start Lumine VPN"

### 3. 创建桌面快捷方式

1. 在 Automate 中，点击流程旁的菜单
2. 选择 "Add to home screen"

**优点**: 
- 可视化编辑
- 更容易理解
- 可以添加更复杂的逻辑

---

## 📋 首次配置 NekoBox

无论使用哪种方案，都需要在 NekoBox 中配置一次 Lumine：

### 手动添加

1. 打开 NekoBox
2. 点击 "+" 按钮
3. 选择 "SOCKS5"
4. 填写:
   - 名称: Lumine
   - 服务器: 127.0.0.1
   - 端口: 1080
5. 保存

### 从剪切板导入

1. 复制: `socks5://127.0.0.1:1080#Lumine`
2. 打开 NekoBox
3. 点击 "从剪切板导入"

---

## 🎯 我的推荐

**最佳方案**: Termux:Widget 桌面快捷方式

**原因**:
- ✓ 完全免费（F-Droid）
- ✓ 一键启动
- ✓ 不需要 root
- ✓ 稳定可靠
- ✓ 容易设置

**设置时间**: 5-10 分钟  
**使用时间**: 1 次点击

---

## ⚙️ 高级配置

### 开机自动启动 (需要 Tasker 或 Automate)

**Tasker**:
1. Profile → Event → System → Device Boot
2. Task: Start Lumine

**Automate**:
1. 添加触发器: Flow Begining → Device Boot

### 监控 Lumine 状态

创建检查脚本 `~/.shortcuts/check-lumine.sh`:

```bash
#!/data/data/com.termux/files/usr/bin/sh

if pgrep -f "liblumine.so" > /dev/null; then
    termux-toast "✓ Lumine 正在运行"
else
    termux-toast "✗ Lumine 未运行"
fi
```

---

## 🔧 故障排除

### 问题: 脚本没有运行

**解决方案**:
```bash
# 给脚本添加执行权限
chmod +x /sdcard/start_lumine_nekobox.sh
# 或
chmod +x ~/.shortcuts/*.sh
```

### 问题: 找不到 liblumine.so

**解决方案**:
1. 确保已安装 NekoBox
2. 检查路径: `ls /data/data/io.nekohasekai.sagernet/lib/liblumine.so`
3. 如果不存在，重新安装 NekoBox

### 问题: 端口被占用

**解决方案**:
```bash
# 查看占用端口 1080 的进程
netstat -tulpn | grep 1080

# 或者修改 Lumine 配置使用其他端口（如 1088）
```

---

## 📞 需要帮助？

如果你在设置过程中遇到问题:
1. 查看本文档的故障排除部分
2. 检查 Termux 的输出信息
3. 使用 `check-lumine.sh` 检查状态

---

## 📚 相关文档

- [LUMINE_QUICK_START.md](../LUMINE_QUICK_START.md) - 快速开始指南
- [LUMINE_USER_GUIDE.md](../LUMINE_USER_GUIDE.md) - 完整用户手册
- [LUMINE_READY_TO_USE_CONFIGS.txt](../LUMINE_READY_TO_USE_CONFIGS.txt) - 即用配置

---

## ✅ 总结

**最简单**: 方案 A (Termux 手动启动)  
**最推荐**: 方案 B (Termux:Widget 桌面快捷方式)  
**最自动**: 方案 C/D (Tasker/Automate 自动化)

选择适合你的方案，享受 Lumine 带来的流量混淆功能！
