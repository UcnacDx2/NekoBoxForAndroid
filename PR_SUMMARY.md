# Lumine 网络预处理插件集成 / Lumine Network Preprocessing Plugin Integration

[中文](#中文说明) | [English](#english-description)

---

## 中文说明

### 概述

本 PR 将 [lumine](https://github.com/UcnacDx2/lumine) 网络预处理工具集成到 NekoBox for Android 中，作为内置插件提供强大的流量混淆和审查绕过功能。

### 主要变更

#### 1. 源代码集成 (3,500+ 行)
- 将 lumine 库移植到 `libcore/lumine/` 目录
- 创建 CLI 包装器 `libcore/lumine_cmd/main.go`
- 添加 libcore 包装器 `libcore/lumine.go`
- 修复所有包声明和依赖关系

#### 2. 构建系统
- 新增 `libcore/build_lumine.sh` - 交叉编译脚本
  - 支持 ARM64、ARMv7、x86、x86_64
  - 输出到 `app/src/main/jniLibs/{arch}/liblumine.so`
- 更新 `libcore/build.sh` 以包含 lumine 构建

#### 3. 插件注册
- `PluginManager.kt`: 添加 "lumine-plugin" 映射
- `PluginEntry.kt`: 添加 Lumine 插件条目
- `Executable.kt`: 添加进程管理支持
- `strings.xml`: 添加显示名称

#### 4. 文档
- `libcore/LUMINE_INTEGRATION.md` - 技术集成文档
- `LUMINE_USER_GUIDE.md` - 双语用户手册
- `libcore/lumine_example_config.json` - 示例配置

#### 5. 安全修复
- 修复 DNS 回退逻辑反转问题 (安全漏洞)
- 修复 HTTP 代理空指针解引用 (崩溃漏洞)
- 修复 SOCKS5 错误处理缺失 (数据损坏风险)

### 功能特性

- ✅ **SOCKS5 代理** (端口 1080)
- ✅ **HTTP 代理** (端口 1081)
- ✅ **TLS 记录分片** - 绕过 DPI 检测
- ✅ **TCP 分段** - 增加检测难度
- ✅ **TTL 脱同步** - 欺骗中间设备
- ✅ **DNS over HTTPS** - 避免 DNS 污染
- ✅ **策略路由** - 基于域名/IP 的灵活规则

### 构建和测试

```bash
# 构建 lumine 插件
cd libcore
./build_lumine.sh

# 构建整个项目 (包括 lumine)
./build.sh
```

### 使用方法

1. **配置文件**: 创建 `/data/data/io.nekohasekai.sagernet/files/lumine_config.json`
2. **启动**: 在 NekoBox 设置中找到 "Lumine" 选项
3. **集成**: 将 Lumine 作为上游代理或在路由规则中使用

详细使用说明请参考 [LUMINE_USER_GUIDE.md](./LUMINE_USER_GUIDE.md)

### 技术架构

Lumine 以**独立进程**运行 (类似 hysteria/naive 插件):
- 进程隔离确保安全性和稳定性
- 独立的生命周期管理
- JSON 配置
- 与 NekoBox 通过本地代理通信

### 文件变更统计

- **修改**: 5 个文件
- **新增**: 22 个文件
- **代码行数**: 3,500+ 行

### 质量保证

- ✅ 代码审查已完成 (发现 6 个问题，修复 3 个关键问题)
- ✅ CodeQL 安全扫描已通过
- ✅ 完整文档 (双语)
- ✅ 示例配置已提供

---

## English Description

### Overview

This PR integrates the [lumine](https://github.com/UcnacDx2/lumine) network preprocessing tool into NekoBox for Android as a built-in plugin, providing powerful traffic obfuscation and censorship circumvention capabilities.

### Major Changes

#### 1. Source Code Integration (3,500+ lines)
- Ported lumine library to `libcore/lumine/` directory
- Created CLI wrapper in `libcore/lumine_cmd/main.go`
- Added libcore wrapper `libcore/lumine.go`
- Fixed all package declarations and dependencies

#### 2. Build System
- New `libcore/build_lumine.sh` - Cross-compilation script
  - Supports ARM64, ARMv7, x86, x86_64
  - Outputs to `app/src/main/jniLibs/{arch}/liblumine.so`
- Updated `libcore/build.sh` to include lumine build

#### 3. Plugin Registration
- `PluginManager.kt`: Added "lumine-plugin" mapping
- `PluginEntry.kt`: Added Lumine plugin entry
- `Executable.kt`: Added process management support
- `strings.xml`: Added display name

#### 4. Documentation
- `libcore/LUMINE_INTEGRATION.md` - Technical integration guide
- `LUMINE_USER_GUIDE.md` - Bilingual user manual
- `libcore/lumine_example_config.json` - Example configuration

#### 5. Security Fixes
- Fixed DNS fallback logic inversion (security issue)
- Fixed HTTP proxy nil pointer dereference (crash vulnerability)
- Fixed missing SOCKS5 error handling (data corruption risk)

### Features

- ✅ **SOCKS5 Proxy** (port 1080)
- ✅ **HTTP Proxy** (port 1081)
- ✅ **TLS Record Fragmentation** - Bypass DPI detection
- ✅ **TCP Segmentation** - Increase detection difficulty
- ✅ **TTL Desynchronization** - Deceive middleboxes
- ✅ **DNS over HTTPS** - Avoid DNS poisoning
- ✅ **Policy Routing** - Flexible domain/IP-based rules

### Build and Test

```bash
# Build lumine plugin
cd libcore
./build_lumine.sh

# Build entire project (including lumine)
./build.sh
```

### Usage

1. **Configuration**: Create `/data/data/io.nekohasekai.sagernet/files/lumine_config.json`
2. **Start**: Find "Lumine" option in NekoBox settings
3. **Integration**: Use Lumine as upstream proxy or in routing rules

For detailed instructions, see [LUMINE_USER_GUIDE.md](./LUMINE_USER_GUIDE.md)

### Technical Architecture

Lumine runs as an **independent process** (like hysteria/naive plugins):
- Process isolation ensures security and stability
- Independent lifecycle management
- JSON configuration
- Communicates with NekoBox via local proxy

### File Change Statistics

- **Modified**: 5 files
- **Added**: 22 files
- **Lines of code**: 3,500+

### Quality Assurance

- ✅ Code review completed (6 issues found, 3 critical fixed)
- ✅ CodeQL security scan passed
- ✅ Complete documentation (bilingual)
- ✅ Example configuration provided

---

## Implementation Details

### Directory Structure

```
libcore/
├── lumine/                    # Lumine library
│   ├── config.go
│   ├── server.go
│   ├── policy.go
│   ├── dns.go
│   ├── fragment.go
│   └── ... (17 Go files total)
├── lumine_cmd/                # CLI wrapper
│   └── main.go
├── lumine.go                  # Libcore wrapper
├── build_lumine.sh            # Build script
└── LUMINE_INTEGRATION.md      # Technical docs

app/
├── src/main/
│   ├── java/.../plugin/
│   │   ├── PluginManager.kt   # Plugin registration
│   │   └── PluginEntry.kt     # Plugin metadata
│   ├── jniLibs/               # Binary output
│   │   ├── arm64-v8a/liblumine.so
│   │   ├── armeabi-v7a/liblumine.so
│   │   ├── x86/liblumine.so
│   │   └── x86_64/liblumine.so
│   └── res/values/strings.xml

LUMINE_USER_GUIDE.md           # User documentation
```

### Key Files Modified

1. **PluginManager.kt** - Added lumine-plugin initialization
2. **PluginEntry.kt** - Added Lumine enum entry
3. **Executable.kt** - Added liblumine.so to execution list
4. **strings.xml** - Added action_lumine string
5. **build.sh** - Added lumine build step

### Commits

1. `455fb83` - Initial plan
2. `a95fcb8` - Integrate lumine network preprocessing tool as Android plugin
3. `7d03308` - Fix critical bugs in lumine code
4. `39d6cc6` - Add Lumine plugin entry and comprehensive user documentation

---

## Credits

- **Original lumine**: https://github.com/UcnacDx2/lumine
- **Based on TlsFragment**: https://github.com/maoist2009/TlsFragment
- **Integration**: NekoBox for Android contributors

## License

GPL-3.0
